package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"kira-go/internal/database"
	"kira-go/internal/models"
)

// geoCache 内存缓存 IP 地理位置，避免重复请求
type geoCache struct {
	mu    sync.RWMutex
	cache map[string]*geoInfo
}

// geoInfo IP 地理位置信息
type geoInfo struct {
	City      string `json:"city"`
	Region    string `json:"region"`
	Country   string `json:"country"`
	District  string `json:"district"`
	Org       string `json:"org"`
	ASN       string `json:"asn"`
	IsMobile  bool   `json:"is_mobile"`
	IsProxy   bool   `json:"is_proxy"`
	IsHosting bool   `json:"is_hosting"`
}

// 全局缓存实例
var geoCacheInstance = &geoCache{
	cache: make(map[string]*geoInfo),
}

// PROVINCE_MAP 英文省份 → 中文 映射
var PROVINCE_MAP = map[string]string{
	"beijing": "北京", "tianjin": "天津", "shanghai": "上海", "chongqing": "重庆",
	"hebei": "河北", "shanxi": "山西", "inner mongolia": "内蒙古",
	"liaoning": "辽宁", "jilin": "吉林", "heilongjiang": "黑龙江",
	"jiangsu": "江苏", "zhejiang": "浙江", "anhui": "安徽", "fujian": "福建",
	"jiangxi": "江西", "shandong": "山东", "henan": "河南", "hubei": "湖北",
	"hunan": "湖南", "guangdong": "广东", "guangxi": "广西", "hainan": "海南",
	"sichuan": "四川", "guizhou": "贵州", "yunnan": "云南", "tibet": "西藏",
	"shaanxi": "陕西", "gansu": "甘肃", "qinghai": "青海", "ningxia": "宁夏",
	"xinjiang": "新疆", "hong kong": "香港", "macau": "澳门", "taiwan": "台湾",
}

// VENDOR_MAP 科技/云厂商 关键词映射
var VENDOR_MAP = map[string]string{
	"alibaba": "阿里巴巴", "aliyun": "阿里云",
	"tencent": "腾讯", "tencent cloud": "腾讯云",
	"baidu": "百度", "baidu cloud": "百度云",
	"huawei": "华为", "huaweicloud": "华为云",
	"bytedance": "字节跳动", "volcengine": "火山引擎",
	"kuaishou": "快手", "meituan": "美团", "xiaomi": "小米",
	"jd.com": "京东", "netease": "网易",
	"kingsoft": "金山云", "qihoo 360": "奇虎 360",
	"sina": "新浪", "sohu": "搜狐", "douban": "豆瓣",
	"amazon": "亚马逊", "aws": "亚马逊云", "microsoft": "微软",
	"azure": "微软云", "google": "谷歌", "google cloud": "谷歌云",
	"cloudflare": "Cloudflare", "oracle": "甲骨文",
	"ovh": "OVH", "digitalocean": "DigitalOcean", "hetzner": "Hetzner",
	"apple": "苹果",
}

// ASN_ORG_MAP 用 asn 号判断运营商
var ASN_ORG_MAP = map[string]string{
	"56041": "中国移动", "56042": "中国移动", "56043": "中国移动",
	"56044": "中国移动", "56045": "中国移动", "56046": "中国移动",
	"56047": "中国移动", "56048": "中国移动", "58453": "中国移动",
	"9808": "中国移动", "24400": "中国移动", "24444": "中国移动", "24445": "中国移动",
	"4808": "中国联通", "4837": "中国联通", "9929": "中国联通", "17816": "中国联通",
	"4134": "中国电信", "4809": "中国电信", "4812": "中国电信",
	"23724": "中国电信", "140061": "中国电信",
	"4538":  "中国教育网",
	"45102": "阿里云", "37963": "阿里云", "45090": "阿里云",
	"45069": "腾讯云", "132203": "腾讯云",
	"136907": "华为云", "55967": "百度云", "137696": "火山引擎",
	"13335": "Cloudflare", "209242": "Cloudflare",
}

// getOrgCN 智能识别运营商/组织名，返回中文
func getOrgCN(org, asn string) string {
	if org == "" {
		return ""
	}

	orgLower := strings.ToLower(org)
	orgLower = strings.ReplaceAll(orgLower, ",", "")
	orgLower = strings.ReplaceAll(orgLower, ".", "")
	orgLower = strings.ReplaceAll(orgLower, "&", "and")

	// 从 org 中提取省份
	var foundProvince string
	for eng, cn := range PROVINCE_MAP {
		if strings.Contains(orgLower, eng) {
			foundProvince = cn
			break
		}
	}

	// 根据 org 中字段 识别运营商大类
	isMobile := strings.Contains(orgLower, "china mobile") || strings.Contains(orgLower, "chinamobile") || strings.Contains(orgLower, "cmcc")
	isUnicom := strings.Contains(orgLower, "china unicom") || strings.Contains(orgLower, "chinaunicom") || strings.Contains(orgLower, "cucc")
	isTelecom := strings.Contains(orgLower, "china telecom") || strings.Contains(orgLower, "chinatelecom") || strings.Contains(orgLower, "chinanet") || strings.Contains(orgLower, "ctcc")
	isCernet := strings.Contains(orgLower, "cernet") || strings.Contains(orgLower, "cernt") || strings.Contains(orgLower, "china education")

	if isMobile {
		if foundProvince != "" {
			return fmt.Sprintf("中国移动 (%s)", foundProvince)
		}
		return "中国移动"
	}
	if isUnicom {
		if foundProvince != "" {
			return fmt.Sprintf("中国联通 (%s)", foundProvince)
		}
		return "中国联通"
	}
	if isTelecom {
		if foundProvince != "" {
			return fmt.Sprintf("中国电信 (%s)", foundProvince)
		}
		return "中国电信"
	}
	if isCernet {
		return "中国教育网"
	}

	// 通用模式识别
	// 含 mobile 但无 china mobile → 地方移动
	if strings.Contains(orgLower, "mobile") || strings.Contains(org, "移动") {
		if foundProvince != "" {
			return fmt.Sprintf("%s 移动", foundProvince)
		}
		return "移动运营商"
	}
	// 含 unicom/united network → 联通系
	if strings.Contains(orgLower, "unicom") || strings.Contains(orgLower, "united network") || strings.Contains(orgLower, "uninet") {
		if foundProvince != "" {
			return fmt.Sprintf("%s 联通", foundProvince)
		}
		return "中国联通"
	}
	// 含 telecom/chinanet → 电信系
	if strings.Contains(orgLower, "telecom") || strings.Contains(orgLower, "chinanet") || strings.Contains(orgLower, "telecommunications") {
		if foundProvince != "" {
			return fmt.Sprintf("%s 电信", foundProvince)
		}
		return "中国电信"
	}
	// 含 netcom → 网通（联通前身）
	if strings.Contains(orgLower, "netcom") || strings.Contains(orgLower, "cnc") {
		return "中国网通"
	}

	// 有线电视 / 广电
	if strings.Contains(orgLower, "cable") {
		if foundProvince != "" {
			return fmt.Sprintf("%s 广电", foundProvince)
		}
		return "有线电视网络"
	}
	if strings.Contains(orgLower, "broadcast") || strings.Contains(orgLower, "radio") || strings.Contains(orgLower, "tv") {
		if foundProvince != "" {
			return fmt.Sprintf("%s 广电", foundProvince)
		}
		return "广电网络"
	}

	// Huashu (华数)
	if strings.Contains(orgLower, "huashu") || strings.Contains(orgLower, "wasu") {
		return "华数传媒"
	}

	// China Networks Inter-Exchange / 互联交换
	if strings.Contains(orgLower, "inter-exchange") || strings.Contains(orgLower, "interexchange") {
		return "互联交换网络"
	}
	if strings.Contains(orgLower, "china networks") {
		return "中国网络交换"
	}

	// 科技/云厂商
	for keyword, cnName := range VENDOR_MAP {
		if strings.Contains(orgLower, keyword) {
			return cnName
		}
	}

	// 从 asn 号反查
	if asn != "" {
		asnNum := strings.ReplaceAll(asn, "AS", "")
		asnNum = strings.ReplaceAll(asnNum, "as", "")
		asnNum = strings.TrimSpace(asnNum)
		if orgName, ok := ASN_ORG_MAP[asnNum]; ok {
			return orgName
		}
	}

	// 从 asn 文本匹配
	if asn != "" {
		asnLower := strings.ToLower(asn)
		isMobileASN := strings.Contains(asnLower, "china mobile") || strings.Contains(asnLower, "chinamobile") || strings.Contains(asnLower, "cmcc")
		isUnicomASN := strings.Contains(asnLower, "china unicom") || strings.Contains(asnLower, "chinaunicom")
		isTelecomASN := strings.Contains(asnLower, "china telecom") || strings.Contains(asnLower, "chinatelecom") || strings.Contains(asnLower, "chinanet")
		if isMobileASN {
			return "中国移动"
		}
		if isUnicomASN {
			return "中国联通"
		}
		if isTelecomASN {
			return "中国电信"
		}
	}

	return ""
}

// uaInfo User-Agent 解析结果
type uaInfo struct {
	Browser    string `json:"browser"`
	OS         string `json:"os"`
	DeviceType string `json:"device_type"`
}

// parseUA 简单的 User-Agent 解析
func parseUA(ua string) uaInfo {
	browser := "Unknown"
	if strings.Contains(ua, "Edg/") {
		browser = "Edge"
	} else if strings.Contains(ua, "Chrome/") {
		browser = "Chrome"
	} else if strings.Contains(ua, "Firefox/") {
		browser = "Firefox"
	} else if strings.Contains(ua, "Safari/") {
		browser = "Safari"
	}

	osName := "Unknown"
	if strings.Contains(ua, "Win") {
		osName = "Windows"
	} else if strings.Contains(ua, "Android") {
		osName = "Android"
	} else if strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad") {
		osName = "iOS"
	} else if strings.Contains(ua, "Mac") {
		osName = "macOS"
	} else if strings.Contains(ua, "Linux") {
		osName = "Linux"
	}

	// 判断设备类型
	deviceType := "电脑"
	if strings.Contains(ua, "Mobi") || strings.Contains(ua, "Android") || strings.Contains(ua, "iPhone") {
		deviceType = "手机"
	}

	return uaInfo{
		Browser:    browser,
		OS:         osName,
		DeviceType: deviceType,
	}
}

// FetchGeo 查询 IP 地理位置（带内存缓存）- 公开函数供 API 层调用
func FetchGeo(ip string) *geoInfo {
	return fetchGeo(ip)
}

// fetchGeo 查询 IP 地理位置（带内存缓存）- 内部函数
func fetchGeo(ip string) *geoInfo {
	// 检查缓存（读锁）
	geoCacheInstance.mu.RLock()
	if cached, ok := geoCacheInstance.cache[ip]; ok {
		geoCacheInstance.mu.RUnlock()
		return cached
	}
	geoCacheInstance.mu.RUnlock()

	// 跳过本地和内网 IP
	if ip == "127.0.0.1" || ip == "::1" || ip == "" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") {
		return nil
	}

	var result *geoInfo

	// 尝试 uapis.cn（国内 API，IPv4/IPv6 支持好，返回中文）
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("https://uapis.cn/api/v1/network/ipinfo?ip=%s", ip))
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var data struct {
				IP     string `json:"ip"`
				Region string `json:"region"`
				ISP    string `json:"isp"`
				ASN    string `json:"asn"`
				As     string `json:"as"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.IP != "" {
				// region 格式："国家 省份 城市" 或 "国家"
				parts := strings.Fields(data.Region)
				result = &geoInfo{
					Country: func() string {
						if len(parts) >= 1 {
							return parts[0]
						}
						return ""
					}(),
					Region: func() string {
						if len(parts) >= 2 {
							return parts[1]
						}
						return ""
					}(),
					City: func() string {
						if len(parts) >= 3 {
							return parts[2]
						}
						return ""
					}(),
					District: "",
					Org:      data.ISP,
					ASN: func() string {
						if data.ASN != "" {
							return data.ASN
						}
						return data.As
					}(),
					IsMobile:  false,
					IsProxy:   false,
					IsHosting: false,
				}
			}
		}
	}

	// 回退：ip-api.com（仅 IPv4）
	if result == nil && !strings.Contains(ip, ":") {
		resp, err = client.Get(fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN&fields=66846719", ip))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var data struct {
					Status     string `json:"status"`
					City       string `json:"city"`
					RegionName string `json:"regionName"`
					Country    string `json:"country"`
					District   string `json:"district"`
					Org        string `json:"org"`
					ISP        string `json:"isp"`
					As         string `json:"as"`
					Mobile     bool   `json:"mobile"`
					Proxy      bool   `json:"proxy"`
					Hosting    bool   `json:"hosting"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.Status == "success" {
					result = &geoInfo{
						City:     data.City,
						Region:   data.RegionName,
						Country:  data.Country,
						District: data.District,
						Org: func() string {
							if data.Org != "" {
								return data.Org
							}
							return data.ISP
						}(),
						ASN:       data.As,
						IsMobile:  data.Mobile,
						IsProxy:   data.Proxy,
						IsHosting: data.Hosting,
					}
				}
			}
		}
	}

	// 写入缓存（写锁）
	if result != nil {
		geoCacheInstance.mu.Lock()
		geoCacheInstance.cache[ip] = result
		geoCacheInstance.mu.Unlock()
	}

	return result
}

// visitorDict 内部数据结构
type visitorDict struct {
	ID         uint   `json:"id"`
	IP         string `json:"ip"`
	Path       string `json:"path"`
	City       string `json:"city"`
	Region     string `json:"region"`
	Country    string `json:"country"`
	District   string `json:"district"`
	Org        string `json:"org"`
	OrgCN      string `json:"org_cn"`
	ASN        string `json:"asn"`
	IsMobile   bool   `json:"is_mobile"`
	IsProxy    bool   `json:"is_proxy"`
	IsHosting  bool   `json:"is_hosting"`
	Browser    string `json:"browser"`
	OS         string `json:"os"`
	DeviceType string `json:"device_type"`
	CreatedAt  string `json:"created_at"`
}

// RecordVisit 记录一次访问
func RecordVisit(ip, path, userAgent string) (*models.Visitor, error) {
	db := database.GetDB()

	// 解析 UA
	uaInfo := parseUA(userAgent)

	// 查询 IP 地理位置
	geoInfo := fetchGeo(ip)

	visitor := &models.Visitor{
		IP:         ip,
		Path:       path,
		UserAgent:  userAgent,
		Browser:    uaInfo.Browser,
		OS:         uaInfo.OS,
		DeviceType: uaInfo.DeviceType,
	}

	if geoInfo != nil {
		visitor.City = geoInfo.City
		visitor.Region = geoInfo.Region
		visitor.Country = geoInfo.Country
		visitor.District = geoInfo.District
		visitor.Org = geoInfo.Org
		visitor.ASN = geoInfo.ASN
		visitor.IsMobile = geoInfo.IsMobile
		visitor.IsProxy = geoInfo.IsProxy
		visitor.IsHosting = geoInfo.IsHosting
	}

	if err := db.Create(visitor).Error; err != nil {
		return nil, err
	}

	return visitor, nil
}

// GetRecentVisitors 获取最近访客列表
func GetRecentVisitors(page, size int) ([]visitorDict, error) {
	db := database.GetDB()

	offset := (page - 1) * size
	var visitors []*models.Visitor
	if err := db.Order("created_at DESC").Offset(offset).Limit(size).Find(&visitors).Error; err != nil {
		return nil, err
	}

	result := make([]visitorDict, 0, len(visitors))
	for _, v := range visitors {
		result = append(result, visitorDict{
			ID:         v.ID,
			IP:         v.IP,
			Path:       v.Path,
			City:       v.City,
			Region:     v.Region,
			Country:    v.Country,
			District:   v.District,
			Org:        v.Org,
			OrgCN:      getOrgCN(v.Org, v.ASN),
			ASN:        v.ASN,
			IsMobile:   v.IsMobile,
			IsProxy:    v.IsProxy,
			IsHosting:  v.IsHosting,
			Browser:    v.Browser,
			OS:         v.OS,
			DeviceType: v.DeviceType,
			CreatedAt:  v.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return result, nil
}

// GetVisitorCount 获取总访客数
func GetVisitorCount() (int64, error) {
	db := database.GetDB()

	var count int64
	if err := db.Model(&models.Visitor{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteVisitor 删除单条访客记录
func DeleteVisitor(visitorID uint) error {
	db := database.GetDB()

	return db.Delete(&models.Visitor{}, visitorID).Error
}

// ClearVisitors 清空所有访客记录
func ClearVisitors() error {
	db := database.GetDB()

	return db.Exec("DELETE FROM visitor").Error
}

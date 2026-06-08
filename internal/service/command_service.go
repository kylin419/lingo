package service

import (
	"fmt"
	"strings"

	"line-translate-bot/internal/repository"
)

var supportedLangs = map[string]string{
	"af":       "Afrikaans",
	"sq":       "Albanian",
	"am":       "Amharic",
	"ar":       "Arabic",
	"hy":       "Armenian",
	"as":       "Assamese",
	"az":       "Azerbaijani (Latin)",
	"bn":       "Bangla",
	"ba":       "Bashkir",
	"eu":       "Basque",
	"bho":      "Bhojpuri",
	"brx":      "Bodo",
	"bs":       "Bosnian (Latin)",
	"bg":       "Bulgarian",
	"yue":      "Cantonese (Traditional)",
	"ca":       "Catalan",
	"hne":      "Chhattisgarhi",
	"lzh":      "Chinese (Literary)",
	"zh-Hans":  "簡體中文",
	"zh-Hant":  "繁體中文",
	"sn":       "ChiShona",
	"hr":       "Croatian",
	"cs":       "Czech",
	"da":       "Danish",
	"prs":      "Dari",
	"dv":       "Divehi",
	"doi":      "Dogri",
	"nl":       "Dutch",
	"en":       "English",
	"et":       "Estonian",
	"fo":       "Faroese",
	"fj":       "Fijian",
	"fil":      "Filipino",
	"fi":       "Finnish",
	"fr":       "French",
	"fr-ca":    "French (Canada)",
	"gl":       "Galician",
	"ka":       "Georgian",
	"de":       "German",
	"el":       "Greek",
	"gu":       "Gujarati",
	"ht":       "Haitian Creole",
	"ha":       "Hausa",
	"he":       "Hebrew",
	"hi":       "Hindi",
	"mww":      "Hmong Daw (Latin)",
	"hu":       "Hungarian",
	"is":       "Icelandic",
	"ig":       "Igbo",
	"id":       "Indonesian",
	"ikt":      "Inuinnaqtun",
	"iu":       "Inuktitut",
	"iu-Latn":  "Inuktitut (Latin)",
	"ga":       "Irish",
	"it":       "Italian",
	"ja":       "Japanese",
	"kn":       "Kannada",
	"ks":       "Kashmiri",
	"kk":       "Kazakh",
	"km":       "Khmer",
	"rw":       "Kinyarwanda",
	"tlh-Latn": "Klingon (Latin)",
	"tlh-Piqd": "Klingon (plqaD)",
	"gom":      "Konkani",
	"ko":       "Korean",
	"ku":       "Kurdish (Central)",
	"kmr":      "Kurdish (Northern)",
	"ky":       "Kyrgyz",
	"lo":       "Lao",
	"lv":       "Latvian",
	"lt":       "Lithuanian",
	"ln":       "Lingala",
	"dsb":      "Lower Sorbian",
	"lug":      "Luganda",
	"mk":       "Macedonian",
	"mai":      "Maithili",
	"mg":       "Malagasy",
	"ms":       "Malay",
	"ml":       "Malayalam",
	"mt":       "Maltese",
	"mni":      "Manipuri",
	"mi":       "Maori",
	"mr":       "Marathi",
	"mn-Cyrl":  "Mongolian (Cyrillic)",
	"mn-Mong":  "Mongolian (Traditional)",
	"my":       "Myanmar (Burmese)",
	"ne":       "Nepali",
	"nb":       "Norwegian Bokmål",
	"ny":       "Nyanja",
	"or":       "Odia",
	"ps":       "Pashto",
	"fa":       "Persian",
	"pl":       "Polish",
	"pt":       "Portuguese (Brazil)",
	"pt-pt":    "Portuguese (Portugal)",
	"pa":       "Punjabi",
	"otq":      "Queretaro Otomi",
	"ro":       "Romanian",
	"run":      "Rundi",
	"ru":       "Russian",
	"sm":       "Samoan",
	"sr-Cyrl":  "Serbian (Cyrillic)",
	"sr-Latn":  "Serbian (Latin)",
	"st":       "Sesotho",
	"nso":      "Sesotho sa Leboa",
	"tn":       "Setswana",
	"sd":       "Sindhi",
	"si":       "Sinhala",
	"sk":       "Slovak",
	"sl":       "Slovenian",
	"so":       "Somali",
	"es":       "Spanish",
	"sw":       "Swahili",
	"sv":       "Swedish",
	"ty":       "Tahitian",
	"ta":       "Tamil",
	"tt":       "Tatar",
	"te":       "Telugu",
	"th":       "Thai",
	"bo":       "Tibetan",
	"ti":       "Tigrinya",
	"to":       "Tongan",
	"tr":       "Turkish",
	"tk":       "Turkmen",
	"uk":       "Ukrainian",
	"hsb":      "Upper Sorbian",
	"ur":       "Urdu",
	"ug":       "Uyghur",
	"uz":       "Uzbek",
	"vi":       "Vietnamese",
	"cy":       "Welsh",
	"xh":       "Xhosa",
	"yo":       "Yoruba",
	"yua":      "Yucatec Maya",
	"zu":       "Zulu",
}

type CommandService struct {
	repo repository.GroupRepository
}

func NewCommandService(repo repository.GroupRepository) *CommandService {
	return &CommandService{repo: repo}
}

func (s *CommandService) HandleSetCommand(sourceID, text string) string {
	input := strings.ReplaceAll(text, ",", " ")
	argsStr := strings.TrimSpace(strings.TrimPrefix(input, "$set"))
	if argsStr == "" {
		return generateHelpMessage("請提供語言代碼！")
	}

	// 使用 strings.Fields 可以自動忽略多餘的空白字元
	langs := strings.Fields(argsStr)
	var validLangs []string
	var invalidLangs []string

	// 檢查使用者輸入的代碼是否存在於我們的字典中
	for _, l := range langs {
		if _, exists := supportedLangs[l]; exists {
			validLangs = append(validLangs, l)
		} else {
			invalidLangs = append(invalidLangs, l)
		}
	}

	// 如果有任何錯誤的代碼，直接拒絕設定並回傳教學
	if len(invalidLangs) > 0 {
		errMsg := fmt.Sprintf("❌ 找不到以下語言代碼：%s", strings.Join(invalidLangs, ", "))
		return generateHelpMessage(errMsg)
	}

	// 儲存合法的設定
	err := s.repo.SaveLanguages(sourceID, validLangs)
	if err != nil {
		return "設定失敗，請稍後再試。"
	}

	// 組合成功提示訊息
	var langNames []string
	for _, l := range validLangs {
		langNames = append(langNames, supportedLangs[l])
	}

	return fmt.Sprintf("✅ 語言已成功設定為：%s\n接下來的訊息我會在這幾個語言間互翻！", strings.Join(langNames, ", "))
}

// generateHelpMessage 負責產生包含可用代碼表的提示訊息
func generateHelpMessage(prefixMsg string) string {
	var helpBuilder strings.Builder
	helpBuilder.WriteString(prefixMsg + "\n\n")
	helpBuilder.WriteString("支援的常用語言代碼如下：\n")

	// 列出所有支援的語言
	for code, name := range supportedLangs {
		helpBuilder.WriteString(fmt.Sprintf("• %s (%s)\n", code, name))
	}
	return helpBuilder.String()
}

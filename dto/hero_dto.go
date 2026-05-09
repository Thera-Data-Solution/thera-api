package dto

type HeroResponse struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	TitleEn       *string `json:"titleEn,omitempty"`
	Subtitle      *string `json:"subtitle,omitempty"`
	SubtitleEn    *string `json:"subtitleEn,omitempty"`
	Description   *string `json:"description,omitempty"`
	DescriptionEn *string `json:"descriptionEn,omitempty"`
	Image         *string `json:"image,omitempty"`
	ButtonText    *string `json:"buttonText,omitempty"`
	ButtonTextEn  *string `json:"buttonTextEn,omitempty"`
	ButtonLink    *string `json:"buttonLink,omitempty"`
	ThemeType     *string `json:"themeType,omitempty"`
}

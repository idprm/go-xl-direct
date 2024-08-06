package entity

type Service struct {
	ID                  int     `json:"id"`
	Category            string  `json:"category"`
	Code                string  `json:"code"`
	Name                string  `json:"name"`
	Package             string  `json:"package"`
	Price               float64 `json:"price"`
	ProgramId           string  `json:"program_id"`
	SidOptIn            string  `json:"sid_optin"`
	SidMt               string  `json:"sid_mt"`
	RenewalDay          int     `json:"renewal_day"`
	TrialDay            int     `json:"trial_day"`
	UrlTelco            string  `json:"url_telco"`
	UrlPortal           string  `json:"url_portal"`
	UrlCallback         string  `json:"url_callback"`
	UrlNotifSub         string  `json:"url_notif_sub"`
	UrlNotifUnsub       string  `json:"url_notif_unsub"`
	UrlNotifRenewal     string  `json:"url_notif_renewal"`
	UrlPostback         string  `json:"url_postback"`
	UrlPostbackBillable string  `json:"url_postback_billable"`
	IsContentSequence   bool    `json:"is_content_sequence"`
}

func (s *Service) GetId() int {
	return s.ID
}

func (s *Service) GetCategory() string {
	return s.Category
}

func (s *Service) GetCode() string {
	return s.Code
}

func (s *Service) GetName() string {
	return s.Name
}

func (s *Service) GetPackage() string {
	return s.Package
}

func (s *Service) GetPrice() float64 {
	return s.Price
}

func (s *Service) GetProgramId() string {
	return s.ProgramId
}

func (s *Service) GetSidOptIn() string {
	return s.SidOptIn
}

func (s *Service) GetSidMt() string {
	return s.SidMt
}

func (s *Service) GetRenewalDay() int {
	return s.RenewalDay
}

func (s *Service) GetTrialDay() int {
	return s.TrialDay
}

func (s *Service) GetUrlTelco() string {
	return s.UrlTelco
}

func (s *Service) GetUrlPortal() string {
	return s.UrlPortal
}

func (s *Service) GetUrlCallback() string {
	return s.UrlCallback
}

func (s *Service) GetUrlNotifSub() string {
	return s.UrlNotifSub
}

func (s *Service) GetUrlNotifUnsub() string {
	return s.UrlNotifUnsub
}

func (s *Service) GetUrlNotifRenewal() string {
	return s.UrlNotifRenewal
}

func (s *Service) GetUrlPostback() string {
	return s.UrlPostback
}

func (s *Service) GetUrlPostbackBillable() string {
	return s.UrlPostbackBillable
}

func (s *Service) GetIsContentSequence() bool {
	return s.IsContentSequence
}

func (s *Service) IsCloudplay() bool {
	return s.GetCategory() == "CLOUDPLAY"
}

func (s *Service) IsGalays() bool {
	return s.GetCategory() == "GALAYS"
}

func (s *Service) IsGupi() bool {
	return s.GetCategory() == "GUPI"
}

func (s *Service) IsMplus() bool {
	return s.GetCategory() == "MPLUS"
}

func (s *Service) IsQuizpro() bool {
	return s.GetCategory() == "QUIZPRO"
}

func (s *Service) IsGameboat() bool {
	return s.GetCategory() == "GAMEBOAT"
}

func (s *Service) IsGamesik() bool {
	return s.GetCategory() == "GAMESIK"
}

func (s *Service) IsDigmagz() bool {
	return s.GetCategory() == "DIGMAGZ"
}

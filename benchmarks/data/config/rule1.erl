rule  "BaseItierRule1"  "ItierRule1"  salience 10 
begin
	if UserFact.RoundNo == 1 && UserFact.Itier != 8 && UserFact.Itier != 12 && UserFact.UserKey =="playone" {
		UserFact.PlayOneAgain = true
	}
	if UserFact.Itier == 20 {
		UserFact.Result="1"
	} else if UserFact.RoundNo == 1 && ( UserFact.Itier == 8 || UserFact.Itier == 12 ) {
		UserFact.Result="2"
	} else if UserFact.RoundNo == 1 && UserFact.Itier != 8 && UserFact.Itier != 12 && UserFact.PlayOneAgain == true {
		UserFact.Result="3"
	} else if UserFact.RoundNo == 1  && UserFact.Itier != 8 && UserFact.Itier != 12 && 
		UserFact.PlayOneAgain == false {
		UserFact.Result="4"
	}else if UserFact.RoundNo <= 5 {
		UserFact.Result="5"
	}
end
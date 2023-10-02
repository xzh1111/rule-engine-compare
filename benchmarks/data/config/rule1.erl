rule  "BaseItierRule1"  "ItierRule1"  salience 10 
begin
	if UserFact.Itier == 20 {
		UserFact.Result="1"
	} else if UserFact.RoundNo == 1 && ( UserFact.Itier == 8 || UserFact.Itier == 12 ) {
		UserFact.Result="2"
	}
end
rule  "GetPlayOneAgain"  "GetPlayOneAgain"  salience 8 
begin
	if UserFact.RoundNo == 1 && UserFact.Itier != 8 && UserFact.Itier != 12 && UserFact.UserKey =="noplayone" {
		UserFact.PlayOneAgain = false
	}
end
rule  "BaseRoundNoRule3"  "RoundNoRule3"  salience 7 
begin
	if 	UserFact.RoundNo == 1 && UserFact.Itier != 8 && UserFact.Itier != 12 && UserFact.PlayOneAgain == true {
		UserFact.Result="2"
	}
end
rule  "BaseRoundNoRule4"  "RoundNoRule4"  salience 6 
begin
	if 	UserFact.RoundNo == 1  && UserFact.Itier != 8 && UserFact.Itier != 12 && 
		UserFact.PlayOneAgain == false {
		UserFact.Result="2"
	}
end
rule  "BaseRoundNoRule5"  "RoundNoRule5"  salience 5 
begin
	if  UserFact.RoundNo <= 5 {
		UserFact.Result="3"
	}
end
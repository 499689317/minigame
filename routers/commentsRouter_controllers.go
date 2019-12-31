package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["minigame/controllers:GameController"] = append(beego.GlobalControllerRouter["minigame/controllers:GameController"],
        beego.ControllerComments{
            Method: "GetGameInfo",
            Router: `/get_game_info`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:GameController"] = append(beego.GlobalControllerRouter["minigame/controllers:GameController"],
        beego.ControllerComments{
            Method: "StartGame",
            Router: `/start_game`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:GameController"] = append(beego.GlobalControllerRouter["minigame/controllers:GameController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "DelCashBonus",
            Router: `/del_cash_bonus`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetCashBonus",
            Router: `/get_cash_bonus`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetCashInfo",
            Router: `/get_cash_info`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetCashList",
            Router: `/get_cash_list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetInviteList",
            Router: `/get_invite_list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "LookTv",
            Router: `/look_tv`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["minigame/controllers:UserController"] = append(beego.GlobalControllerRouter["minigame/controllers:UserController"],
        beego.ControllerComments{
            Method: "ShareGame",
            Router: `/share_game`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

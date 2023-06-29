package handler

type RequestGuacamole struct {
	GuacamoleAddr string `json:"guacamole_addr"`
	AssetProtocol string `json:"asset_protocol"`
	AssetHost     string `json:"asset_host"`
	AssetPort     string `json:"asset_port"`
	AssetUser     string `json:"asset_user"`
	AssetPassword string `json:"asset_password"`
	ScreenWidth   int    `json:"screen_width"`
	ScreenHeight  int    `json:"screen_height"`
	ScreenDpi     int    `json:"screen_dpi"`
}

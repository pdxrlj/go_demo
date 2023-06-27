package handler

type RequestGuacamole struct {
	GuacamoleAddr string
	AssetProtocol string
	AssetHost     string
	AssetPort     string
	AssetUser     string
	AssetPassword string
	ScreenWidth   int
	ScreenHeight  int
	ScreenDpi     int
}

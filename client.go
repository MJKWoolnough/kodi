package kodi

type Client interface {
	ExecuteAddon(name string, params ...interface{}) string

	Quit()
	Mute(bool)
	SetVolume(uint8) // 0 -> 100
}

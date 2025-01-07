# vlctom3u
A tool to export .m3u8 playlists from the VLC andoid app.

# Usage
1. Dump the vlc media database from the app (Settings -> Advanced -> Dump media database)
2. Clone the repository and place the vlc_media.db file in the same directory
3. `go run main.go`
4. The exported .m3u8 playlists will appear in the same directory
5. Enjoy!

# todo
This tool does not support internet streams and other playlist entries that don't have a folder_id

Feel free to contribute to my shitty code :P

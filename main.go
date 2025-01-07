package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "net/url"
    _ "github.com/mattn/go-sqlite3"
)

type Playlist struct {
    id   string
    name string
}

type Song struct {
    artist string
    title string
    media_id string
    position string
    filename string
    folder_id string
    path string
    duration string
}

func main() {
    db, err := sql.Open("sqlite3", "./vlc_media.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT id_playlist, name FROM Playlist")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var playlists []Playlist
    for rows.Next() {
        var playlist Playlist
        err := rows.Scan(&playlist.id, &playlist.name)
        if err != nil {
            log.Fatal(err)
        }
        playlists = append(playlists, playlist)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Opened db, exporting playlists to m3u")
    for _, playlist := range playlists {
        fmt.Printf("------------------%s------------------\n", playlist.name)
    
        rows, err = db.Query("SELECT position, media_id FROM PlaylistMediaRelation WHERE playlist_id = ?", playlist.id)

        var songs []Song
        for rows.Next() {
            var song Song
            err := rows.Scan(&song.position, &song.media_id)
            if err != nil {
                log.Fatal(err)
            }
            songs = append(songs, song)
        }

        for i, song := range songs {
            var vlcduration string
            var artist_id string
            err = db.QueryRow("SELECT title, filename, folder_id, duration, artist_id FROM Media WHERE id_media = ?", song.media_id).Scan(&song.title, &song.filename, &song.folder_id, &vlcduration, &artist_id)
            song.duration = fmt.Sprintf("%.3s", vlcduration)
            if err != nil {
                log.Print(err)
                continue
            }

            var urlpath string
            err = db.QueryRow("SELECT path FROM Folder WHERE id_folder = ?", song.folder_id).Scan(&urlpath)
            if err != nil {
                log.Print(err)
                continue
            }

            err = db.QueryRow("SELECT name FROM Artist WHERE id_artist = ?", artist_id).Scan(&song.artist)
            if err != nil {
                log.Print(err)
            }

            path, err := url.PathUnescape(urlpath)
            if err != nil {
                log.Print(err)
            }

            song.path = path + song.filename
            songs[i] = song
        }

        file, err := os.Create(playlist.name + ".m3u8")
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        for _, song := range songs {
            fmt.Println(song.title)
            _, err = file.WriteString(song.path + "\n")
            if err != nil {
                log.Fatal(err)
            }
        }

        file.Close()
    }
    fmt.Println("")
    fmt.Println("Done :)")
}
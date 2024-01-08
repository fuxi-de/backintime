// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func spotifyWebplayer() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_spotifyWebplayer_18bc`,
		Function: `function __templ_spotifyWebplayer_18bc(){console.log('initting spotify player')
  const token = localStorage.getItem("backInTime-token")
  console.log(token)
  window.onSpotifyWebPlaybackSDKReady = () => {
        console.log('onSpotifyWebPlaybackSDKReady')
        const player = new window.Spotify.Player({
            name: 'Web Playback SDK',
            getOAuthToken: cb => { cb(token); },
        });
        player.addListener('ready', ({ device_id }) => {
            console.log('Ready with Device ID', device_id);
            window.backintime = {} 
            window.backintime.device_id = device_id
            window.backintime.player = player
        });

        player.addListener('not_ready', ({ device_id }) => {
            console.log('Device ID has gone offline', device_id);
        });

        player.connect();

    };}`,
		Call:       templ.SafeScript(`__templ_spotifyWebplayer_18bc`),
		CallInline: templ.SafeScriptInline(`__templ_spotifyWebplayer_18bc`),
	}
}

func SpotifyPlayer() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = spotifyWebplayer().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script async src=\"https://sdk.scdn.co/spotify-player.js\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := ``
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

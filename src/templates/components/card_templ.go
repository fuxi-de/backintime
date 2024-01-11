// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func triggerPlayback(category string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_triggerPlayback_5ae6`,
		Function: `function __templ_triggerPlayback_5ae6(category){const token = localStorage.getItem("backInTime-token")
  const playButton = document.querySelector(".active .play")
  const revealButton = document.querySelector(".active .reveal")
  console.log(window.backintime)
    playButton.addEventListener("click", () => {
      fetch('http://localhost:1312/user/play/'+window.backintime.device_id+'/'+category, {
		    method: 'GET',
		    headers: {
			    'Authorization': 'Bearer '+token,
		    },
	    })
      playButton.classList.add("hidden")
      revealButton.classList.remove("hidden")
    })}`,
		Call:       templ.SafeScript(`__templ_triggerPlayback_5ae6`, category),
		CallInline: templ.SafeScriptInline(`__templ_triggerPlayback_5ae6`, category),
	}
}

func getSongInfo(category string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_getSongInfo_81aa`,
		Function: `function __templ_getSongInfo_81aa(category){const token = localStorage.getItem("backInTime-token")
    const revealButton = document.querySelector(".active .reveal")
    const activeCard = document.querySelector(".active")
    const addButton = document.querySelector(".add-card-button")
    revealButton.addEventListener("click", () => {
      window.backintime.player.getCurrentState().then( state => { 
        if (!state) {
          console.error('User is not playing music through the Web Playback SDK');
          return;
        }

        htmx.ajax('GET', ` + "`" + `/user/play/release/${state.track_window.current_track.id}` + "`" + `, { target: ".active .year", headers: { Authorization: "Bearer "+token}})

        var current_track = state.track_window.current_track;
        var next_track = state.track_window.next_tracks[0];

        window.backintime.player.pause()
        const wrapper = document.querySelector(".active")
        const title = wrapper.querySelector(".title")
        title.textContent = current_track.name
        const interpret = wrapper.querySelector(".interpret")
        interpret.textContent = current_track.artists[0].name
        addButton.disabled = false
        activeCard.classList.remove('active')
      });
    })}`,
		Call:       templ.SafeScript(`__templ_getSongInfo_81aa`, category),
		CallInline: templ.SafeScriptInline(`__templ_getSongInfo_81aa`, category),
	}
}

func Card(category string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div style=\"box-shadow:8px 8px black\" class=\"w-fit border-4 border-solid border-black p-3 active\"><div class=\"group relative\"><img class=\"w-full md:w-72 block rounded opacity-0 \" src=\"https://upload.wikimedia.org/wikipedia/en/f/f1/Tycho_-_Epoch.jpg\" alt=\"\"><div class=\"absolute rounded w-full h-full top-0 flex items-center justify-evenly bg-transparent group-hover:bg-black group-hover:cursor-pointer\"><button class=\"play text-black scale-150 group-hover:text-white\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"40\" height=\"40\" fill=\"currentColor\" class=\"bi bi-play-circle-fill\" viewBox=\"0 0 16 16\"><path d=\"M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zM6.79 5.093A.5.5 0 0 0 6 5.5v5a.5.5 0 0 0 .79.407l3.5-2.5a.5.5 0 0 0 0-.814l-3.5-2.5z\"></path></svg></button> <button class=\"reveal hidden group-hover:text-white\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := `reveal`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button></div></div><div class=\"p-5\"><h3 class=\"title text-black text-lg\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var3 := `Title`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h3><p class=\"interpret text-gray-700\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var4 := `Interpret`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var4)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><p class=\"year text-gray-700\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var5 := `Year`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var5)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = triggerPlayback(category).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = getSongInfo(category).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

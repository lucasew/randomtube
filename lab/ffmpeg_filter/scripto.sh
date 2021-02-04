rm out.mp4 || true > /dev/null
export FILTER='
[0:v]scale=1920*2:-1,crop=1920:1080,boxblur=luma_radius=min(h\,w)/20:luma_power=1:chroma_radius=min(cw\,ch)/20:chroma_power=1[bg];
[0:v]scale=-1:1080[ov];
[bg][ov]overlay=(W-w)/2,crop=1920:1080
'

echo $FILTER
ffmpeg -y -i src.mp4 -filter_complex "$FILTER" -c:a copy out.mp4 && vlc out.mp4 --fullscreen --play-and-exit

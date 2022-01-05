rm out.mp4 || true > /dev/null
export FILTER='
[0:v]scale=1920*2:-1,crop=1920:1080,boxblur=luma_radius=min(h\,w)/20:luma_power=1:chroma_radius=min(cw\,ch)/20:chroma_power=1[bg];
[0:v]scale=-1:1080[ov];
[bg][ov]overlay=(W-w)/2,crop=1920:1080[composed];
anullsrc=cl=stereo:r=44100[anull]
'

echo $FILTER
ffmpeg -y  -i src.mp4 -t 15 -filter_complex "$FILTER" -map '[composed]' -map '0:a?' -map '[anull]' -shortest out.mkv && vlc out.mkv --fullscreen --play-and-exit

async (event, steps, params, auths) => {
    const axios = require('axios');
    const authorized_chat = params.authorized_chat
    const {host} = steps.trigger.event.headers
    const url = new URL(event.url)
    async function justReturn() {
        return await $respond({
            status: 200,
            body: "ok"
        })
    }
    const that = this;
    if (!$checkpoint) {
        $checkpoint = {}
    }
    if (!$checkpoint.videos) {
        $checkpoint.videos = []
    }
    async function invokeTelegram(method, data) {
        return await axios({
            method: "GET",
            url: `https://api.telegram.org/bot${auths.telegram_bot_api.token}/${method}`,
            data
        })
    }
    const handler = {
        async deleteIds() {
            const isTelegramAllowed = !event.query.disallowTelegram
            console.log(isTelegramAllowed)
            const isTelegram = event.headers["user-agent"].indexOf("Telegram") !== -1

            if ((isTelegram && isTelegramAllowed) || !isTelegram) {
                const bodyFileIds = event.body.files !== undefined 
                    ? event.body.files.length > 0 ? event.body.files : []
                    : [];
                const queryFileIds = event.query.ids != undefined
                    ? event.query.ids.split(',')
                    : [];
                const fileIds = [...bodyFileIds, ...queryFileIds]
                $checkpoint.videos = $checkpoint.videos.filter((video) => {
                    return fileIds.indexOf(video.file_id) === -1
                })
                console.log(fileIds)
            } else {
                console.log("Telegram requested and not allowed")
            }
            await justReturn()
        },
        async notify() {
            await invokeTelegram("sendMessage", {
                text: event.body.messsage || `mensagem vazia enviada por ${event.client_ip}`,
                chat_id: authorized_chat
            })
            return justReturn()
        },
        async yta() {
            return await $respond({
                immediate: true,
                status: 200,
                body: auths.youtube_data_api
            })
        },
        async telegram() {
            async function pushVideo({file_id, duration, file_size}) {
                if (file_size > (20*1024*1024)) {
                    return await reply(`erro: ${file_size} bytes é mais que 20MB, não vou conseguir baixar`)
                }
                $checkpoint.videos.push({
                    file_id: file_id,
                    length: duration
                })
                await reply(`tá na lista, patrão\nremover: ${host}/deleteIds?ids=${file_id}&disallowTelegram=1`)
            }
            async function reply(text, parse_mode) {
                await invokeTelegram("sendMessage", {
                    text, 
                    chat_id: event.body.message.chat.id,
                    reply_to_message_id: event.body.message.message_id,
                    parse_mode,
                })
            }
            if (event.body.message.chat.id !== authorized_chat) {
                await justReturn()
                return await invokeTelegram("sendMessage", {
                    text: `Usuário @${event.body.message.from.username} tentou usar o bot sem autorização pelo chat ${event.body.message.chat.username}`,
                    chat_id: authorized_chat
                })
            }
            const {video, animation, message_id, document} = event.body.message
            if (video) {
                await pushVideo(video)
            } else if (animation) {
                await pushVideo(animation)
            } else if (document && document.mime_type === "video/mp4") {
                await reply("quase lá, manda para o @convert2filebot que ele resolve esse arquivo")
            } else {
                await reply("erro: não é vídeo nem gif")
            }    
            return await justReturn()
        },
        async list() {
            return await $respond({
                immediate: true,
                status: 200,
                body: {
                    ...$checkpoint,
                    reportChatID: params.authorized_chat
                }
            })
        },
        async len() {
            let secs = 0
            const videos = $checkpoint.videos
            let qt = videos.length
            for (let i = 0; i < qt; i++) {
                secs += videos[i].length || 0
            }
            return await $respond({
                immediate: true,
                status: 200,
                body: {
                    secs,
                    qt
                }
            })
        },
        "favicon.ico": async function() {
            return await $respond({
                immediate: true,
                status: 404
            })
        }
    }[url.pathname.slice(1)]

    handler && await handler()
}

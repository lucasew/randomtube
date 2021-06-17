# RandomTube

A system to make Youtube compilation videos from selected ones on a Telegram chat.

This works using two components:

- The pipedream agent
  - Reacts to Telegram messages using Telegram webhooks.
    - Checks if a message comes from a trusted chat (in our case a group chat).
    - Checks if Telegram is requesting the endpoint from the preview system or for a real link using the `User-Agent`.
  - Leaks the Youtube token to be used on the other component.
  - Store Telegram `file_id`s to be downloaded later
    - It also checks if the item is a video and doesn't have more than 20MB
- The hard work agent
  - Does the video processing
    - Requires FFMpeg
    - Normalize different video formats into a common one
      - This is where the real interesting part is xD
  - Fetches the YouTube and Telegram token from Pipedream
  - Joins the normalized video
  - Post the video on YouTube
  - It can run as a cron job and on GitHub Actions

Both agents work together. The pipedream agent is always running, at least when it's awaken, the hard work agent can be run manually and run sometimes.
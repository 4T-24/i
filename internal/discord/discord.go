package discord

import (
	"instancer/internal/env"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

var (
	client webhook.Client
)

func InitDiscordClient() {
	cfg := env.Get()
	client = webhook.New(snowflake.ID(cfg.Discord.WebhookId), cfg.Discord.WebhookToken)
}

func SendMessage(title, message string, color int) error {
	_, err := client.CreateMessage(discord.NewWebhookMessageCreateBuilder().
		SetUsername("Instancer").
		SetEmbeds(discord.NewEmbedBuilder().
			SetTitle(title).
			SetDescription(message).
			SetColor(color).
			Build(),
		).
		Build())
	return err
}

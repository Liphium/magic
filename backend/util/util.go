package util

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"os"

	"github.com/Liphium/magic/backend/util/constants"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// This is just useless and there just for fun, this is justified performance wasting
var quotes []string = []string{
	"It's cool to have you along fellow travellor.",
	"None of this is AI-generated, so yes, I wrote all of these just by hand.",
	"Did you start building yet?",
	"Did you know you can even test Minecraft plugins on here?",
	"Liphium used to be just a chat app. Now we've even got our own build pipelines, please help!",
	"All of this isn't magic, but it sometimes feels like the others are using debuff magic.",
	"Any ideas for AI features for this platform? We need some investors.",
	"Welcome to your cloud-native, blazingly fast, memory-safe and bleeding-edge cloud platform.",
	"\"Why do you still test like this when testing can feel like Magic?\"",
	"Real software engineers test in production, but only when it's not actually production.",
	"Random quotes are cool, add them to your app!",
	"We need to make these quotes change daily and not be random, that would be cool!",
	"Welcome back, Builder 💪! Don't read random text, you need to hussle!",
	"If we ever have some sort of feature showcase, it should be called \"Magic Show\"!",
	"The Magic logo could've had two more crosses in it. So it would look like an M..",
	"Magic isn't just in Isekai, it's right here in front of you!",
	"The urge to write \"Legal jargon\" instead of \"Legal documents\", it almost happened.",
	"Welcome back to Magic! Have you figured out how to create fire yet?",
	"Have you been to the Magic Forge? I hear they make the best swords there.",
	"How many quotes do you think there are? If you guess correctly, we'll take you seriously.",
	"Rumors are saying we actually created fire magic, where are they coming from?",
	"Water, fire, air or ice magic: Which one is your favorite?",
	"It's like magic ✨✨",
	"Magic isn't actually open-source, but I hope that's fine.",
	"These quotes are almost like shower thoughts, you never know what's coming next.",
	"Have you had dinner today? Cause I'm about to get dinner after finishing these quotes.",
	"Knock knock. Who's there? Square. Square who? DeepSeek won't know.",
	"I feel like a certain guy should stop gambling..",
	"There is only one correct ranking of programming languages: Go > all.",
	"Let's Go to the Magic Show together!",
	"Magic is magical like magic in a magic show.",
	"Why is this called Magic? Because these quotes magically appear ✨✨",
}

func RandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateToken(tkLength int32) string {

	s := make([]rune, tkLength)

	length := big.NewInt(int64(len(letters)))

	for i := range s {

		number, _ := crand.Int(crand.Reader, length)
		s[i] = letters[number.Int64()]
	}

	return string(s)
}

// Get the account uuid
func AccountUUID(c *fiber.Ctx) uuid.UUID {
	uuid, _ := uuid.Parse(c.Locals(constants.LocalsAccountID).(string))
	return uuid
}

// Check if the backend server is in testing mode
func IsTesting() bool {
	return os.Getenv("MAGIC_TESTING") == "true"
}

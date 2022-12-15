package config

import "github.com/matorix/chat.fdbm.dev/internal/utils"

const SecretKey = "3k69xkjcreuK#"
const BodyLengthLimit = 10000
const NameLengthLimit = 14
const DatabaseListLimit = 50
const DatabaseName = "database.db"

var BaseDir, BaseDirErr = utils.GetBaseDir()

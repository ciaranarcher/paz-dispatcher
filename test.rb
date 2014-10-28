#!/usr/bin/env ruby

require "redis"
require "json"

redis = Redis.new(:url => "redis://127.0.0.1")

message_id = (rand() * 65535).to_i
ssid = (rand() * 15).to_i
content = {
  "callsign" => "PY2CHM-#{ssid}",
  "subject" => "Hey this is a subject #{message_id}",
  "description" => "Whats'up? Whatching the game, having a bud. True True"
}
puts content.to_json
redis.rpush("paz:inq", content.to_json)

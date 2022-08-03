package main

const CONFIG_FILE_PATH = "./config.json"

const DEFAULT_DELIMITTER = "^"
const DEFAULT_OUTPUT_DIRECTORY = "./out"
const DEFAULT_LAYERS_DIRECTORY = "./layers"

var INVALID_DELIMITTERS = []string{".", ",", ""}

const DEFAULT_GIF_DISPOSAL uint8 = 1
const DEFAULT_GIF_DELAY = 50

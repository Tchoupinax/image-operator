import pino, { LoggerOptions } from "pino";

const config: LoggerOptions = {
  level: "DEBUG",
  base: null,
};

if (process.env.NODE_ENV !== "production") {
  config.transport = {
    target: "pino-pretty",
    options: {
      colorize: true,
      levelFirst: true,
      translateTime: "HH:MM:ss.l",
    },
  };
}

export const logger = pino(config);

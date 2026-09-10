declare module "slack" {
  export type JSONValue = null | boolean | number | string | JSONValue[] | {[key: string]: JSONValue};
  export interface MessageRef {channelId: string; ts: string}
  export interface Text {text: string}
  export interface Context {
    id: string; teamId: string; channelId: string; userId: string;
    command: string; text: string;
    event: {type: string; text: string; ts: string; threadTs: string; channelId: string; userId: string};
    config: Record<string, string | boolean | number>;
    reply(message: Text): Promise<MessageRef | {delivered: true; via: "response_url"}>;
    slack: {messages: {post(message: Text & {channelId: string; threadTs?: string}): Promise<MessageRef>}};
    store: {get(key: string): JSONValue | undefined; set(key: string, value: JSONValue): void; delete(key: string): boolean; keys(): string[]};
    log: {debug(message: string): void; info(message: string): void; warn(message: string): void; error(message: string): void};
  }
  export interface Registration {
    configure(spec: {name: string; description?: string; run?: {fields: Record<string, {type: "string" | "bool" | "number"; default?: string | boolean | number; required?: boolean; help?: string}>}}): void;
    command(name: string, spec: {description: string}, handler: (ctx: Context) => Text | void | Promise<Text | void>): void;
    event(name: "app_mention", handler: (ctx: Context) => Text | void | Promise<Text | void>): void;
  }
  export function defineBot(register: (api: Registration) => void): object;
}

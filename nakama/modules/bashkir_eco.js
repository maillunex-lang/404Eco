/**
 * Nakama JS runtime module — Bashkortostan Economic Simulator
 * RPCs bridge multiplayer auth/presence with the Go game-api.
 */

const GAME_API = (typeof process !== "undefined" && process.env && process.env.GAME_API_URL)
  ? process.env.GAME_API_URL
  : "http://game-api:8080";

function rpcRegister(context, logger, nk, payload) {
  let body = {};
  try {
    body = JSON.parse(payload || "{}");
  } catch (e) {
    throw Error("invalid json");
  }
  const name = body.name || context.username || "Предприниматель";
  const res = nk.httpRequest(
    GAME_API + "/api/register",
    "post",
    { "Content-Type": "application/json" },
    JSON.stringify({ nakamaId: context.userId, name: name })
  );
  if (res.code < 200 || res.code >= 300) {
    logger.error("register failed: %s", res.body);
    throw Error(res.body || "register failed");
  }
  return res.body;
}

function rpcSnapshot(context, logger, nk, payload) {
  const res = nk.httpRequest(GAME_API + "/api/snapshot", "get", {}, null);
  if (res.code < 200 || res.code >= 300) {
    throw Error("snapshot failed");
  }
  return res.body;
}

function rpcCreateFirm(context, logger, nk, payload) {
  let body = {};
  try {
    body = JSON.parse(payload || "{}");
  } catch (e) {
    throw Error("invalid json");
  }
  // Ensure player exists
  rpcRegister(context, logger, nk, JSON.stringify({ name: body.playerName || context.username }));
  const reg = JSON.parse(
    nk.httpRequest(
      GAME_API + "/api/register",
      "post",
      { "Content-Type": "application/json" },
      JSON.stringify({ nakamaId: context.userId, name: body.playerName || context.username })
    ).body
  );
  body.playerId = reg.id;
  const res = nk.httpRequest(
    GAME_API + "/api/firms",
    "post",
    { "Content-Type": "application/json" },
    JSON.stringify(body)
  );
  if (res.code < 200 || res.code >= 300) {
    throw Error(res.body || "create firm failed");
  }
  // Broadcast to match channel
  try {
    nk.streamSend(
      { mode: 2, label: "bashkir-world" },
      JSON.stringify({ type: "firm_created", firm: JSON.parse(res.body) }),
      true
    );
  } catch (e) {
    logger.warn("stream send: %s", e.message);
  }
  return res.body;
}

function rpcFirmAction(context, logger, nk, payload) {
  let body = {};
  try {
    body = JSON.parse(payload || "{}");
  } catch (e) {
    throw Error("invalid json");
  }
  const reg = JSON.parse(
    nk.httpRequest(
      GAME_API + "/api/register",
      "post",
      { "Content-Type": "application/json" },
      JSON.stringify({ nakamaId: context.userId, name: context.username })
    ).body
  );
  body.playerId = reg.id;
  const res = nk.httpRequest(
    GAME_API + "/api/firms/action",
    "post",
    { "Content-Type": "application/json" },
    JSON.stringify(body)
  );
  if (res.code < 200 || res.code >= 300) {
    throw Error(res.body || "firm action failed");
  }
  return res.body;
}

function rpcPlaceOrder(context, logger, nk, payload) {
  let body = {};
  try {
    body = JSON.parse(payload || "{}");
  } catch (e) {
    throw Error("invalid json");
  }
  const reg = JSON.parse(
    nk.httpRequest(
      GAME_API + "/api/register",
      "post",
      { "Content-Type": "application/json" },
      JSON.stringify({ nakamaId: context.userId, name: context.username })
    ).body
  );
  body.playerId = reg.id;
  const res = nk.httpRequest(
    GAME_API + "/api/orders",
    "post",
    { "Content-Type": "application/json" },
    JSON.stringify(body)
  );
  if (res.code < 200 || res.code >= 300) {
    throw Error(res.body || "order failed");
  }
  return res.body;
}

function afterAuthenticateDevice(ctx, logger, nk, data) {
  logger.info("device auth user=%s", ctx.userId);
  try {
    nk.httpRequest(
      GAME_API + "/api/register",
      "post",
      { "Content-Type": "application/json" },
      JSON.stringify({ nakamaId: ctx.userId, name: ctx.username || "Игрок" })
    );
  } catch (e) {
    logger.warn("auto-register: %s", e.message);
  }
}

// Match handler for realtime world lobby
function matchInit(context, logger, nk, params) {
  return {
    state: {
      tick: 0,
      presences: {},
      label: "Республика Башкортостан",
    },
    tickRate: 1,
    label: "bashkir-lobby",
  };
}

function matchJoinAttempt(context, logger, nk, dispatcher, tick, state, presence, metadata) {
  return { state: state, accept: true };
}

function matchJoin(context, logger, nk, dispatcher, tick, state, presences) {
  for (var i = 0; i < presences.length; i++) {
    state.presences[presences[i].userId] = presences[i];
  }
  dispatcher.broadcastMessage(1, JSON.stringify({
    type: "welcome",
    message: "Добро пожаловать в экономику Башкортостана",
    online: Object.keys(state.presences).length,
  }));
  return { state: state };
}

function matchLeave(context, logger, nk, dispatcher, tick, state, presences) {
  for (var i = 0; i < presences.length; i++) {
    delete state.presences[presences[i].userId];
  }
  return { state: state };
}

function matchLoop(context, logger, nk, dispatcher, tick, state, messages) {
  state.tick = tick;
  if (tick % 5 === 0) {
    try {
      var res = nk.httpRequest(GAME_API + "/api/nakama/sync", "post", { "Content-Type": "application/json" }, "{}");
      if (res.code >= 200 && res.code < 300) {
        dispatcher.broadcastMessage(2, res.body);
      }
    } catch (e) {
      // game-api may be restarting
    }
  }
  for (var i = 0; i < messages.length; i++) {
    dispatcher.broadcastMessage(3, messages[i].data, null, messages[i].sender);
  }
  return { state: state };
}

function matchTerminate(context, logger, nk, dispatcher, tick, state, graceSeconds) {
  return { state: state };
}

function matchSignal(context, logger, nk, dispatcher, tick, state, data) {
  return { state: state, data: data };
}

// Init
function InitModule(ctx, logger, nk, initializer) {
  initializer.registerRpc("eco_register", rpcRegister);
  initializer.registerRpc("eco_snapshot", rpcSnapshot);
  initializer.registerRpc("eco_create_firm", rpcCreateFirm);
  initializer.registerRpc("eco_firm_action", rpcFirmAction);
  initializer.registerRpc("eco_place_order", rpcPlaceOrder);
  initializer.registerAfterAuthenticateDevice(afterAuthenticateDevice);
  initializer.registerMatch("bashkir_lobby", {
    matchInit: matchInit,
    matchJoinAttempt: matchJoinAttempt,
    matchJoin: matchJoin,
    matchLeave: matchLeave,
    matchLoop: matchLoop,
    matchTerminate: matchTerminate,
    matchSignal: matchSignal,
  });
  logger.info("Bashkortostan Eco Nakama module loaded");
}

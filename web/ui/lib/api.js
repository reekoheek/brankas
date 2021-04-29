import { EventEmitter } from 'events';

const DEFAULT_BASE = typeof location === 'object' ? `${location.protocol}//${location.host}` : 'http://localhost';
const DEFAULT_FETCH = typeof fetch === 'function' ? fetch : undefined;
const DEFAULT_HEADERS = {
  'Content-Type': 'application/json',
};

export class API extends EventEmitter {
  constructor ({ base = DEFAULT_BASE, fetch = DEFAULT_FETCH } = {}) {
    super();

    this.base = base;
    this.fetch = fetch;
    this.headers = {};
  }

  get logined () {
    return !!this.token;
  }

  setToken (token) {
    if (token) {
      this.user = parseJwt(token);
      this.token = token;
      this.headers.Authorization = `Bearer ${token}`;
    } else {
      delete this.user;
      delete this.token;
      delete this.headers.Authorization;
    }

    this.emit('change');
  }

  async req ({ uri, method = 'GET', body }) {
    const fetch = this.fetch;
    const resp = await fetch(`${this.base}${uri}`, {
      method,
      headers: {
        ...DEFAULT_HEADERS,
        ...this.headers,
      },
      body: JSON.stringify(body),
    });

    const respBody = await resp.json();

    if (resp.status >= 200 && resp.status < 400) {
      return respBody;
    }

    throw new HTTPError({
      status: resp.status,
      message: respBody.message,
      children: respBody.errors,
    });
  }

  get (uri) {
    return this.req({ uri });
  }

  post (uri, body) {
    return this.req({ uri, method: 'POST', body });
  }

  put (uri, body) {
    return this.req({ uri, method: 'PUT', body });
  }

  delete (uri) {
    return this.req({ uri, method: 'DELETE' });
  }

  async login ({ username, password }) {
    const { token } = await this.post('/auth/login', { username, password });
    this.setToken(token);
  }

  logout () {
    this.setToken();
  }
}

class HTTPError extends Error {
  constructor ({ status, field = '', message = '', children = [] }) {
    super(message);

    this.status = status;
    this.field = field;
    this.children = children.map(err => new HTTPError({
      field: err.field,
      message: err.message,
    }));
  }
}

const atobFn = typeof atob === 'undefined' ? require('atob') : atob;
function parseJwt (token) {
  const base64Url = token.split('.')[1];
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  const jsonPayload = decodeURIComponent(atobFn(base64).split('').map(c => {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));

  return JSON.parse(jsonPayload);
}

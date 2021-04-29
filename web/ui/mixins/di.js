export function Requester (Base) {
  return class extends Base {
    requestInstance (key) {
      const evt = new CustomEvent('request-instance', {
        bubbles: true,
        detail: { key },
      });

      this.dispatchEvent(evt);

      return evt.detail.instance;
    }
  };
}

export function Provider (Base) {
  return class extends Base {
    constructor () {
      super();

      this._instances = {};
    }

    connectedCallback () {
      super.connectedCallback();

      this._requestInstanceCallback = evt => {
        evt.detail.instance = this._instances[evt.detail.key];
      };

      this.addEventListener('request-instance', this._requestInstanceCallback, false);
    }

    disconnectedCallback () {
      super.disconnectedCallback();

      this.removeEventListener('request-instance', this._requestInstanceCallback, false);
    }

    provideInstance (key, instance) {
      this._instances[key] = instance;
    }

    requestInstance (key) {
      return this._instances[key];
    }
  };
}

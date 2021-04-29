import { LitElement, html } from 'lit-element';
import { Provider } from './mixins/di';
// import { API } from '../lib/api';

// import 'litx-router';

// const LOADERS = [
//   {
//     test () { return true; },
//     load (uri) { return import(`./views/${uri}`); },
//   },
// ];

class XApp extends Provider(LitElement) {
  // static get properties () {
  //   return {
  //     loginUser: { type: String },
  //   };
  // }

  // constructor () {
  //   super();

  //   this.api = new API();
  //   this.api.addListener('change', () => {
  //     if (this.api.logined) {
  //       this.loginUser = this.api.user.sub;
  //       localStorage.API_TOKEN = this.api.token;
  //     } else {
  //       this.loginUser = '';
  //       delete localStorage.API_TOKEN;
  //     }
  //   });

  //   if (localStorage.API_TOKEN) {
  //     this.api.setToken(localStorage.API_TOKEN);
  //   }

  //   this.provideInstance('api', this.api);
  // }

  render () {
    return html`
      <nav class="navbar navbar-light bg-light">
        <div class="container-fluid">
          <a class="navbar-brand" href="/">Brankas</a>
          <button class="btn btn-primary" @click="${this.showMenuOrLogin}">
            <i class="${this.loginUser ? 'bi-unlock' : 'bi-lock'}"></i>
            ${this.loginUser}
          </button>
        </div>
      </nav>
    `;
  }

  createRenderRoot () {
    return this;
  }

  connectedCallback () {
    super.connectedCallback();

    document.body.removeAttribute('unresolved');

    // setTimeout(() => {
    //   const router = this.querySelector('#router');
    //   // router.use(async (ctx, next) => {
    //   //   await next();
    //   // });
    //   this.provideInstance('router', router);

    //   router.start();
    // });
  }
}

customElements.define('x-app', XApp);

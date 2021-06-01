import { FunctionComponent } from 'react';
import { Form } from './components';
import './scss/styles.scss';

interface Props {}

const App: FunctionComponent<Props> = () => {
  const params = new URLSearchParams(window.location.search);

  const clientId = params.get('client_id');
  const state = params.get('state');
  const nonce = params.get('nonce');
  const redirectUri = params.get('redirect_uri');
  const responseType = params.get("response_type");
  const scope = params.get("scope");

  return (
    <main className="form-login">
      <h1 className="mb-3 display-4">Login</h1>

      <Form
        clientId={clientId}
        state={state}
        nonce={nonce}
        redirectUri={redirectUri}
        responseType={responseType}
        scope={scope}
      />

      <p className="mt-5 mb-3 text-muted">
        Plant Pot &copy; {new Date().getFullYear()}
      </p>
    </main>
  );
};

export default App;

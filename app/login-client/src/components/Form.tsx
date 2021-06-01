import {
  ChangeEventHandler,
  FormEventHandler,
  FunctionComponent,
  useState,
} from 'react';
import { login } from '../api';
import { ErrorModel, LoginModel, LoginResponseModel } from '../models';
import FormView from './FormView';

export interface FormProps {
  clientId: string | null;
  state: string | null;
  nonce: string | null;
  redirectUri: string | null;
  responseType: string | null;
  scope: string | null;
}

const Form: FunctionComponent<FormProps> = ({
  state,
  nonce,
  clientId,
  redirectUri,
  responseType,
  scope,
}) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<ErrorModel | null>(null);

  const handleSubmit: FormEventHandler<HTMLFormElement> = async e => {
    e.preventDefault();

    if (loading) {
      return;
    }

    setLoading(true);

    const data: LoginModel = {
      email,
      password,
      state,
      nonce,
      clientId,
      redirectUri,
      responseType,
      scopes: scope?.split(" ") ?? [],
    };

    const res = await login(data);
    if ((res as ErrorModel)?.error) {
      setError(res as ErrorModel);
    } else {
      const data = res as LoginResponseModel;
      window.location.replace(data.redirectUri);
    }

    setLoading(false);
  };

  const handleChange: ChangeEventHandler<HTMLInputElement> = e => {
    const { name, value } = e.target;

    switch (name) {
      case 'email':
        setEmail(value);
        break;
      case 'password':
        setPassword(value);
        break;
    }
  };

  return (
    <FormView
      email={email}
      password={password}
      loading={loading}
      error={error}
      onSubmit={handleSubmit}
      onChange={handleChange}
    />
  );
};

export default Form;

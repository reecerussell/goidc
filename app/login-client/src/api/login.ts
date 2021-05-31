import { ErrorModel, LoginModel, LoginResponseModel } from '../models';

const loginUrl = (): string => {
  const match = window.location.href.match(/prod|dev|test/);
  if (!match || match.length < 1){
    return "/oauth/authorize"
  }

  const stage = match[0];
  return `/${stage}/oauth/authorize`
}

const login = async (
  data: LoginModel
): Promise<ErrorModel | LoginResponseModel> => {
  const res = await fetch(loginUrl(), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (res.status !== 200) {
    if (res.headers.get('Content-Type')?.indexOf('application/json') === -1) {
      return {
        error: 'An error occurred while communicating with the server.',
      };
    }

    return (await res.json()) as ErrorModel;
  }

  return (await res.json()) as LoginResponseModel;
};

export { login };

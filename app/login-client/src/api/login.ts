import { ErrorModel, LoginModel, LoginResponseModel } from '../models';

const loginUrl =
  process.env.BASE_URL + process.env.NODE_ENV === 'production'
    ? '/prod/oauth/login'
    : '/dev/oauth/login';

const login = async (
  data: LoginModel
): Promise<ErrorModel | LoginResponseModel> => {
  const res = await fetch(loginUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (res.status !== 200) {
    if (res.headers.get('Content-Type')?.indexOf('application/json') === -1) {
      return {
        type: res.statusText,
        statusCode: 500,
        message: 'An error occurred while communicating with the server.',
      };
    }

    return (await res.json()) as ErrorModel;
  }

  return (await res.json()) as LoginResponseModel;
};

export { login };

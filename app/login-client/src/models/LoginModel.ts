export default interface LoginModel {
  clientId: string | null;
  state: string | null;
  nonce: string | null;
  redirectUri: string | null;
  responseType: string | null;
  scopes: string[];
  email: string;
  password: string;
}

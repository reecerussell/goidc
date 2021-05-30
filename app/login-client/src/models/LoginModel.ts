export default interface LoginModel {
  clientId: number;
  state: string | null;
  nonce: string | null;
  redirectUri: string | null;
  responseType: string | null;
  scope: string | null;
  email: string;
  password: string;
}

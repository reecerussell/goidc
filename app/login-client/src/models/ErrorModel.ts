export default interface ErrorModel {
  type: string;
  message: string;
  statusCode: number;
  paramName?: string;
}

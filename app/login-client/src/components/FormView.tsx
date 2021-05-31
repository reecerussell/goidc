import { ChangeEventHandler, FormEventHandler, FunctionComponent } from 'react';
import { ErrorModel } from '../models';

export interface FormViewProps {
  email: string;
  password: string;
  loading: boolean;
  error: ErrorModel | null;
  onSubmit: FormEventHandler<HTMLFormElement>;
  onChange: ChangeEventHandler<HTMLInputElement>;
}

const FormView: FunctionComponent<FormViewProps> = ({
  email,
  password,
  loading,
  error,
  onSubmit,
  onChange,
}) => (
  <>
    {error &&  (
      <div className="alert alert-danger" role="alert">
        {error.error}
      </div>
    )}

    <form onSubmit={onSubmit}>
      <div className="form-floating">
        <input
          type="email"
          autoComplete="username"
          className="form-control"
          required
          id="email"
          name="email"
          placeholder="name@example.com"
          value={email}
          onChange={onChange}
        />
        <label htmlFor="email">Email</label>
      </div>
      <div className="form-floating">
        <input
          type="password"
          autoComplete="password"
          className="form-control"
          required
          id="password"
          name="password"
          placeholder="Password"
          value={password}
          onChange={onChange}
        />
        <label htmlFor="password">Password</label>
      </div>

      <button
        type="submit"
        className="w-100 btn btn-lg btn-primary"
        disabled={loading}
      >
        {loading ? (
          <>
            <span
              className="spinner-grow spinner-grow-sm"
              role="status"
              aria-hidden="true"
            ></span>
            <span className="visually-hidden">Loading...</span>
          </>
        ) : (
          'Login'
        )}
      </button>
    </form>
  </>
);

export default FormView;

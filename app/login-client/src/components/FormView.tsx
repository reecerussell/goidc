import { ChangeEventHandler, FormEventHandler, FunctionComponent } from 'react';
import classNames from 'classnames';
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
    {error && error.paramName === null && (
      <div className="alert alert-danger" role="alert">
        {error.message}
      </div>
    )}

    <form onSubmit={onSubmit}>
      <div className="form-floating">
        <input
          type="email"
          autoComplete="username"
          required
          className={classNames('form-control', {
            'is-invalid': error?.paramName === 'email',
          })}
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
          required
          className={classNames('form-control', {
            'is-invalid': error?.paramName === 'password',
          })}
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

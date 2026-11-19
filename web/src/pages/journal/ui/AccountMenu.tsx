interface AccountMenuProps {
  email: string;
  onLogout: () => void;
  pending: boolean;
}

// Первая буква email в рукописном кружке; email и выход — в нативном popover
// (Esc и клик мимо закрывают его сами).
export function AccountMenu({ email, onLogout, pending }: AccountMenuProps) {
  return (
    <>
      <button type="button" className="account__toggle" popoverTarget="account-menu" aria-label="аккаунт">
        <svg className="account__ring" viewBox="0 0 32 32" aria-hidden="true">
          <path d="M15.6 3.6c6.9-.5 12.6 4.5 12.6 11.6 0 7-5.6 12.6-12.6 12.3C9 27.2 4 22 4.1 15.6 4.3 9.3 9.2 4 15.6 3.6z" />
        </svg>
        <span className="account__initial" aria-hidden="true">
          {email.charAt(0)}
        </span>
      </button>
      <div id="account-menu" className="account__menu" popover="auto">
        <p className="account__email">{email}</p>
        <button type="button" className="account__logout" onClick={onLogout} disabled={pending}>
          Выйти
        </button>
      </div>
    </>
  );
}

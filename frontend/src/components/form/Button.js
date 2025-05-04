import styles from './Button.module.css';

const Button = ({type = "button", value, onClick = () => {}}) => {
  return (
    <button
      className={styles.button}
      type={type}
      onClick={onClick}
    >
      {value}
    </button>
  );
};

export default Button;

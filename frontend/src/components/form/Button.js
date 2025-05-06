import styles from './Button.module.css';

const Button = ({type = "button", value, onClick = () => {}, style, className}) => {
  return (
    <button
      className={`${styles.button} ${className || ''}`}
      type={type}
      onClick={onClick}
      style={style}
    >
      {value}
    </button>
  );
};

export default Button;

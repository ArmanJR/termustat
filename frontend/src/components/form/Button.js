import styles from './Button.module.css';

const Button = ({type, value}) => {
  return (
    <button className={styles.button} type={type}>{value}</button>
  );
};

export default Button;

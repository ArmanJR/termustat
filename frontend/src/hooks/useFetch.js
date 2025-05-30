import { useState, useEffect } from "react";

const useFetch = (fetchFunction, deps = null) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetch = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetchFunction();
      setData(response.data);
    } catch (error) {
      if (error.response?.status === 500) {
        setError("مشکلی در سرور رخ داده است. لطفا دوباره تلاش کنید.");
      } else {
        setError(error.message);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetch();
  }, deps === null ? [] : deps);

  return { data, loading, error, fetchData: fetch };
};

export default useFetch;

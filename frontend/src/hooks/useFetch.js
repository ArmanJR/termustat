import { useState, useEffect } from "react";

const useFetch = (fetchFunction, id = null) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetch = async () => {
    setLoading(true);
    setError(null);

    try {
      let response;
      if (id) {
        response = await fetchFunction(id);
      } else {
        response = await fetchFunction();
      }
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
  }, []);

  return { data, loading, error, fetchData: fetch };
};

export default useFetch;

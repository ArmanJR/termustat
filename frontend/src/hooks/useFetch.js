import { useState, useEffect, useRef } from "react";

const useFetch = (fetchFunction, id = null) => {
  const hasFetched = useRef(false);
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
    if (hasFetched.current) return;
    hasFetched.current = true;
    fetch();
  }, [fetchFunction, id]);

  return { data, loading, error, fetchData: fetch };
};

export default useFetch;

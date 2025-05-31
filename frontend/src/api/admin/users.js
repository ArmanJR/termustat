import axiosInstance from "../axiosInstance";

export const getUsers = async (page = 1) => {
  const response = await axiosInstance.get(
    `/admin/users`,
    {
      params: { page },
    }
  );
  return response;
};

export const deleteUser = async (id) => {
  const response = await axiosInstance.delete(
    `/admin/users/${id}`
  );
  return response;
};

export const editUser = async (id, userData) => {
  const response = await axiosInstance.put(
    `/admin/users/${id}`,
    userData
  );
  return response;
};

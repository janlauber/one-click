// a function which validates a string to be in a valid format for a dns subdomain
// e.g. "test" is valid, "test." is not valid, "test-asdf" is valid, "test_asdf" is not valid, truncate to 63 characters
export const isValidName = (str: string): boolean => {
    return /^[a-z0-9-]{1,63}$/.test(str);
};

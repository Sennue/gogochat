-- psql -f schema.sql
-- \i schema.sql

-- Switch Role and Database
\c gogochat
SET ROLE gogochat;
SET search_path TO gogochat,public;

-- Data Generator
CREATE OR REPLACE FUNCTION xdy(rolls integer, sides integer)
	RETURNS integer AS
$$
DECLARE
	result integer := 0;
BEGIN
	FOR i IN 1..rolls LOOP
		result = result + trunc(1 + sides*random());
	END LOOP;
	RETURN result;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION randomString(strings varchar[])
	RETURNS varchar AS
$$
BEGIN
	RETURN strings[xdy(1, array_length(strings, 1))];
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateName(length integer)
	RETURNS varchar AS
$$
DECLARE
	result varchar = '';
	consonant varchar[] = ARRAY[	'b', 'c', 'd', 'f', 'g',
					'h', 'j', 'k', 'l', 'm',
					'n', 'p', 'qu', 'r', 's',
					't', 'v', 'w', 'x', 'z'];
	vowel varchar[] = ARRAY[	'a', 'e', 'i', 'o', 'u', 'y'];
	ofs integer = xdy(1, 2);
BEGIN
	FOR i IN 1..length LOOP
		IF 0 = (ofs + i) % 2 THEN
			result = result || randomString(consonant);
		ELSE
			result = result || randomString(vowel);
		END IF;
	END LOOP;
	RETURN initcap(result);
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateFullName()
	RETURNS varchar AS
$$
DECLARE
	result varchar = '';
	middle_initials integer = xdy(2,2)-2;
BEGIN
	result = result || generateName(xdy(2,4));
	FOR i IN 1..middle_initials LOOP
		result = result || ' ' || left(generateName(1), 1);
	END LOOP;
	result = result || ' ' || generateName(xdy(2,6));
	RETURN result;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateEmail(username varchar, hostlength integer)
	RETURNS varchar AS
$$
DECLARE
	result varchar = '';
BEGIN
	result = result || lower(replace(username, ' ', '.')) || '@';
	result = result || generateName(hostlength) || '.com';
	RETURN lower(result);
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateDeviceId()
	RETURNS varchar AS
$$
DECLARE
	result varchar = '';
BEGIN
	result = result || left(generateName(8), 8) || '-';
	result = result || left(generateName(4), 4) || '-';
	result = result || left(generateName(4), 4) || '-';
	result = result || left(generateName(4), 4) || '-';
	result = result || left(generateName(12), 12);
	RETURN 'FakeID:' || upper(result);
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateUser(population bigint)
	RETURNS TABLE(name varchar, email varchar, password varchar, salt varchar) AS
$$
BEGIN
	FOR i IN 1..population LOOP
		name := generateFullName();
		email := generateEmail(name, 1+xdy(2,2));
		-- http://stackoverflow.com/questions/2647158/how-can-i-hash-passwords-in-postgresql
		salt := gen_salt('bf', 8);
		password := crypt('password', salt);
		RETURN NEXT;
	END LOOP;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateDevice(population bigint)
	RETURNS TABLE(id varchar, account_id bigint) AS
$$
BEGIN
	FOR i IN 1..population LOOP
		id := generateDeviceId();
		account_id := i;
		RETURN NEXT;
	END LOOP;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateWord(length integer)
	RETURNS varchar AS
$$
BEGIN
	return lower(generateName(length));
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateSentence(words integer)
	RETURNS varchar AS
$$
DECLARE
	result varchar = '';
BEGIN
	result = result || initcap(generateWord(xdy(2,4)));
	FOR i IN 2..words LOOP
		result = result || ' ' || generateWord(xdy(2,4));
	END LOOP;
	result = result || '.  ';
	RETURN result;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION generateDeterministicData(population bigint)
	RETURNS VOID AS
$$
DECLARE
	oldseed double precision = random();
BEGIN
	PERFORM setseed(pi() / 10);
	INSERT INTO account (name, email, password, salt)
		SELECT * FROM generateUser(population);
	INSERT INTO device (device_id, account_id)
		SELECT generateDeviceId(), account_id FROM account;
	FOR i IN 1..population LOOP
		INSERT INTO room (name, description)
			SELECT initcap(generateSentence(xdy(1,3))), generateSentence(xdy(2,4));
	END LOOP;
	FOR i IN 1..3 LOOP
		INSERT INTO message (account_id, room_id, body)
			SELECT account_id, room_id, generateSentence(xdy(1,9)) FROM account CROSS JOIN room;
	END LOOP;
	PERFORM setseed(oldseed);
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION deleteData()
	RETURNS VOID AS
$$
BEGIN
	DELETE FROM message;
	DELETE FROM room;
	DELETE FROM device;
	DELETE FROM account;
END;
$$ LANGUAGE 'plpgsql';

-- Initial Data
SELECT deleteData();
SELECT generateDeterministicData(5);

-- Data Test
SELECT * from account;
SELECT * from device;
SELECT * from room;
SELECT * from message;
\dt gogochat.*


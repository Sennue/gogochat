-- psql -f schema.sql
-- \i schema.sql

-- Data Generator
CREATE OR REPLACE FUNCTION gogochat.xdy(rolls integer, sides integer)
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

CREATE OR REPLACE FUNCTION gogochat.randomString(strings varchar[])
	RETURNS varchar AS
$$
BEGIN
	RETURN strings[gogochat.xdy(1, array_length(strings, 1))];
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION gogochat.generateName(length integer)
	RETURNS varchar AS
$$
DECLARE
	result varchar = '';
	consonant varchar[] = ARRAY[	'b', 'c', 'd', 'f', 'g',
					'h', 'j', 'k', 'l', 'm',
					'n', 'p', 'qu', 'r', 's',
					't', 'v', 'w', 'x', 'z'];
	vowel varchar[] = ARRAY[	'a', 'e', 'i', 'o', 'u', 'y'];
	ofs integer = gogochat.xdy(1, 2);
BEGIN
	FOR i IN 1..length LOOP
		IF 0 = (ofs + i) % 2 THEN
			result = result || gogochat.randomString(consonant);
		ELSE
			result = result || gogochat.randomString(vowel);
		END IF;
	END LOOP;
	RETURN initcap(result);
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION gogochat.generateFullName()
	RETURNS varchar AS
$$
DECLARE
	result varchar = '';
	middle_initials integer = gogochat.xdy(2,2)-2;
BEGIN
	result = result || gogochat.generateName(gogochat.xdy(2,4));
	FOR i IN 1..middle_initials LOOP
		result = result || ' ' || left(gogochat.generateName(1), 1);
	END LOOP;
	result = result || ' ' || gogochat.generateName(gogochat.xdy(2,6));
	RETURN result;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION gogochat.generateUser(population integer)
	RETURNS TABLE(name varchar, age int) AS
$$
BEGIN
	FOR i IN 1..population LOOP
		name := gogochat.generateFullName();
		age :=  gogochat.xdy(3,25);
		RETURN NEXT;
	END LOOP;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION gogochat.generateDeterministicUser(population integer)
	RETURNS TABLE(name varchar, age int) AS
$$
DECLARE
	oldseed double precision = random();
BEGIN
	PERFORM setseed(pi() / 10);
	RETURN QUERY SELECT * FROM gogochat.generateUser(population);
	PERFORM setseed(oldseed);
END;
$$ LANGUAGE 'plpgsql';

-- Initial Data
INSERT INTO gogochat.user (name, age)
	SELECT * FROM gogochat.generateDeterministicUser(22);

-- Data Test
SELECT * from gogochat.user;
SELECT COUNT(*) from gogochat.user;
\dt gogochat.*


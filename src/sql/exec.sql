create table item(id integer primary key autoincrement
                 , name text unique
                 , trans text
                 , descr text);


CREATE TABLE cases (id integer primary key autoincrement
				   , id_item integer references item(id)
				   , use_case text unique
				   , comment text
);
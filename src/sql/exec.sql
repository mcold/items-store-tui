create table item(id integer primary key autoincrement
                 , item text unique
                 , trans text
                 , descr text);


CREATE TABLE cases (id integer primary key autoincrement
				   , id_item integer references item(id)
				   , use_case text
				   , comment text
);
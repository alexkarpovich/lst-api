CREATE EXTENSION IF NOT EXISTS ltree;

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TABLE IF EXISTS languages;
CREATE TABLE languages (
  id serial PRIMARY KEY,
  code VARCHAR(2) UNIQUE NOT NULL,
  iso_name VARCHAR(64) NOT NULL,
  native_name VARCHAR(64) NOT NULL
);

DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id serial PRIMARY KEY,
  email VARCHAR(128) UNIQUE NOT NULL,
  username VARCHAR(128) UNIQUE NOT NULL,
  encrypted_password VARCHAR(128) NOT NULL,
  first_name VARCHAR (64),
  last_name VARCHAR (64),
  token VARCHAR (128),
  status SMALLINT NOT NULL,
  token_expires_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER updated_at BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

DROP TABLE IF EXISTS groups;
CREATE TABLE groups (
  id serial PRIMARY KEY,
  target_lang_id INT NOT NULL,
  native_lang_id INT NOT NULL,
  name VARCHAR(128) UNIQUE NOT NULL,
  status SMALLINT NOT NULL,
  CONSTRAINT fk_target_lang
    FOREIGN KEY(target_lang_id) 
    REFERENCES languages(id),
  CONSTRAINT fk_native_lang
    FOREIGN KEY(native_lang_id) 
    REFERENCES languages(id)
);

DROP TABLE IF EXISTS user_group;
CREATE TABLE user_group (
  user_id INT NOT NULL,
  group_id INT NOT NULL,
  role SMALLINT NOT NULL,
  CONSTRAINT fk_user
    FOREIGN KEY(user_id) 
    REFERENCES users(id),
  CONSTRAINT fk_group
    FOREIGN KEY(group_id) 
    REFERENCES groups(id)
);

DROP TABLE IF EXISTS slices;
CREATE TABLE slices (
  id serial PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  visibility SMALLINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TRIGGER updated_at BEFORE UPDATE ON slices
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


DROP TABLE IF EXISTS object_types;
CREATE TABLE object_types (
  id serial PRIMARY KEY,
  name VARCHAR(32) UNIQUE NOT NULL
);


DROP TABLE IF EXISTS expressions;
CREATE TABLE expressions (
  id serial PRIMARY KEY,
  lang_id INT NOT NULL,
  value VARCHAR(128) NOT NULL,
  CONSTRAINT fk_lang
    FOREIGN KEY(lang_id) 
    REFERENCES languages(id),
  UNIQUE(lang_id, value)
);

DROP TABLE IF EXISTS translations;
CREATE TABLE translations (
  id serial PRIMARY KEY,
  target_id INT NOT NULL,
  native_id INT NOT NULL,
  type INT NOT NULL,
  comment VARCHAR(128),
  CONSTRAINT fk_type
    FOREIGN KEY(type) 
    REFERENCES object_types(id),
  UNIQUE(type, target_id, native_id)
);

DROP TABLE IF EXISTS texts;
CREATE TABLE texts (
  id serial PRIMARY KEY,
  origin_id INT,
  lang_id INT NOT NULL,
  author_id INT NOT NULL,
  title VARCHAR(128) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_origin
    FOREIGN KEY(origin_id) 
    REFERENCES texts(id),
  CONSTRAINT fk_lang
    FOREIGN KEY(lang_id) 
    REFERENCES languages(id),
  CONSTRAINT fk_author
    FOREIGN KEY(author_id) 
    REFERENCES users(id)
);

DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
  id serial PRIMARY KEY,
  parent_id INT,
  author_id INT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_author
    FOREIGN KEY(author_id) 
    REFERENCES users(id),
  CONSTRAINT fk_parent
    FOREIGN KEY(parent_id) 
    REFERENCES comments(id)
);
CREATE TRIGGER updated_at BEFORE UPDATE ON comments
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


DROP TABLE IF EXISTS group_slice;
CREATE TABLE group_slice (
  group_id INT NOT NULL,
  slice_id INT NOT NULL,
  path ltree NOT NULL,
  CONSTRAINT fk_group
    FOREIGN KEY(group_id) 
    REFERENCES groups(id),
  CONSTRAINT fk_slice
    FOREIGN KEY(slice_id) 
    REFERENCES slices(id)
);

CREATE INDEX path_gist_idx ON group_slice USING GIST (path);
CREATE INDEX path_idx ON group_slice USING BTREE (path);

DROP TABLE IF EXISTS slice_expression;
CREATE TABLE slice_expression (
  slice_id INT NOT NULL,
  expression_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_slice
    FOREIGN KEY(slice_id) 
    REFERENCES slices(id),
  CONSTRAINT fk_expression
    FOREIGN KEY(expression_id) 
    REFERENCES expressions(id),
  UNIQUE(slice_id, expression_id)
);

DROP TABLE IF EXISTS slice_translation;
CREATE TABLE slice_translation (
  slice_id INT NOT NULL,
  translation_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_slice
    FOREIGN KEY(slice_id) 
    REFERENCES slices(id),
  CONSTRAINT fk_translation
    FOREIGN KEY(translation_id) 
    REFERENCES translations(id),
  UNIQUE(slice_id, translation_id)
);

DROP TABLE IF EXISTS slice_text;
CREATE TABLE slice_text (
  slice_id INT NOT NULL,
  text_id INT NOT NULL,
  CONSTRAINT fk_slice
    FOREIGN KEY(slice_id) 
    REFERENCES slices(id),
  CONSTRAINT fk_text
    FOREIGN KEY(text_id) 
    REFERENCES texts(id),
  UNIQUE(slice_id, text_id)
);

DROP TABLE IF EXISTS object_comment;
CREATE TABLE object_comment (
  object_id INT NOT NULL,
  comment_id INT NOT NULL,
  type INT NOT NULL,
  CONSTRAINT fk_comment
    FOREIGN KEY(comment_id) 
    REFERENCES comments(id),
  CONSTRAINT fk_type
    FOREIGN KEY(type) 
    REFERENCES object_types(id),
  UNIQUE(type, comment_id, object_id)
);

/* INITIALIZE LANGUAGES */

INSERT INTO languages (code, iso_name, native_name) VALUES
('ab', 'Abkhaz', 'аҧсуа'),
('aa', 'Afar', 'Afaraf'),
('af', 'Afrikaans', 'Afrikaans'),
('ak', 'Akan', 'Akan'),
('sq', 'Albanian', 'Shqip'),
('am', 'Amharic', 'አማርኛ'),
('ar', 'Arabic', 'العربية'),
('an', 'Aragonese', 'Aragonés'),
('hy', 'Armenian', 'Հայերեն'),
('as', 'Assamese', 'অসমীয়া'),
('av', 'Avaric', 'авар мацӀ'),
('ae', 'Avestan', 'avesta'),
('ay', 'Aymara', 'aymar aru'),
('az', 'Azerbaijani', 'azərbaycan dili'),
('bm', 'Bambara', 'bamanankan'),
('ba', 'Bashkir', 'башҡорт теле'),
('eu', 'Basque', 'euskara, euskera'),
('be', 'Belarusian', 'Беларуская'),
('bn', 'Bengali', 'বাংলা'),
('bh', 'Bihari', 'भोजपुरी'),
('bi', 'Bislama', 'Bislama'),
('bs', 'Bosnian', 'bosanski jezik'),
('br', 'Breton', 'brezhoneg'),
('bg', 'Bulgarian', 'български език'),
('my', 'Burmese', 'ဗမာစာ'),
('ca', 'Catalan; Valencian', 'Català'),
('ch', 'Chamorro', 'Chamoru'),
('ce', 'Chechen', 'нохчийн мотт'),
('ny', 'Chichewa', 'chiCheŵa'),
('zh', 'Chinese', '中文'),
('cv', 'Chuvash', 'чӑваш чӗлхи'),
('kw', 'Cornish', 'Kernewek'),
('co', 'Corsican', 'corsu'),
('cr', 'Cree', 'ᓀᐦᐃᔭᐍᐏᐣ'),
('hr', 'Croatian', 'hrvatski'),
('cs', 'Czech', 'česky'),
('da', 'Danish', 'dansk'),
('dv', 'Divehi', 'ދިވެހި'),
('nl', 'Dutch', 'Nederlands'),
('en', 'English', 'English'),
('eo', 'Esperanto', 'Esperanto'),
('et', 'Estonian', 'eesti'),
('ee', 'Ewe', 'Eʋegbe'),
('fo', 'Faroese', 'føroyskt'),
('fj', 'Fijian', 'vosa Vakaviti'),
('fi', 'Finnish', 'suomi'),
('fr', 'French', 'français'),
('ff', 'Fula', 'Fulfulde'),
('gl', 'Galician', 'Galego'),
('ka', 'Georgian', 'ქართული'),
('de', 'German', 'Deutsch'),
('el', 'Greek', 'Ελληνικά'),
('gn', 'Guaraní', 'Avañeẽ'),
('gu', 'Gujarati', 'ગુજરાતી'),
('ht', 'Haitian', 'Kreyòl ayisyen'),
('ha', 'Hausa', 'هَوُسَ'),
('he', 'Hebrew', 'עברית'),
('iw', 'Hebrew', 'עברית'),
('hz', 'Herero', 'Otjiherero'),
('hi', 'Hindi', 'हिन्दी'),
('ho', 'Hiri Motu', 'Hiri Motu'),
('hu', 'Hungarian', 'Magyar'),
('ia', 'Interlingua', 'Interlingua'),
('id', 'Indonesian', 'Bahasa Indonesia'),
('ie', 'Interlingue', 'Originally called Occidental'),
('ga', 'Irish', 'Gaeilge'),
('ig', 'Igbo', 'Asụsụ Igbo'),
('ik', 'Inupiaq', 'Iñupiaq'),
('io', 'Ido', 'Ido'),
('is', 'Icelandic', 'Íslenska'),
('it', 'Italian', 'Italiano'),
('iu', 'Inuktitut', 'ᐃᓄᒃᑎᑐᑦ'),
('ja', 'Japanese', '日本語'),
('jv', 'Javanese', 'basa Jawa'),
('kl', 'Kalaallisut', 'kalaallisut'),
('kn', 'Kannada', 'ಕನ್ನಡ'),
('kr', 'Kanuri', 'Kanuri'),
('ks', 'Kashmiri', 'कश्मीरी'),
('kk', 'Kazakh', 'Қазақ тілі'),
('km', 'Khmer', 'ភាសាខ្មែរ'),
('ki', 'Kikuyu', 'Gĩkũyũ'),
('rw', 'Kinyarwanda', 'Ikinyarwanda'),
('ky', 'Kirghiz', 'кыргыз тили'),
('kv', 'Komi', 'коми кыв'),
('kg', 'Kongo', 'KiKongo'),
('ko', 'Korean', '한국어'),
('ku', 'Kurdish', 'Kurdî'),
('kj', 'Kwanyama', 'Kuanyama'),
('la', 'Latin', 'latine'),
('lb', 'Luxembourgish', 'Lëtzebuergesch'),
('lg', 'Luganda', 'Luganda'),
('li', 'Limburgish', 'Limburgs'),
('ln', 'Lingala', 'Lingála'),
('lo', 'Lao', 'ພາສາລາວ'),
('lt', 'Lithuanian', 'lietuvių kalba'),
('lu', 'Luba-Katanga', 'Luba-Katanga'),
('lv', 'Latvian', 'latviešu valoda'),
('gv', 'Manx', 'Gaelg'),
('mk', 'Macedonian', 'македонски јазик'),
('mg', 'Malagasy', 'Malagasy fiteny'),
('ms', 'Malay', 'بهاس ملايو'),
('ml', 'Malayalam', 'മലയാളം'),
('mt', 'Maltese', 'Malti'),
('mi', 'Māori', 'te reo Māori'),
('mr', 'Marathi', 'मराठी'),
('mh', 'Marshallese', 'Kajin M̧ajeļ'),
('mn', 'Mongolian', 'монгол'),
('na', 'Nauru', 'Ekakairũ Naoero'),
('nv', 'Navajo', 'Diné bizaad'),
('nb', 'Norwegian Bokmål', 'Norsk bokmål'),
('nd', 'North Ndebele', 'isiNdebele'),
('ne', 'Nepali', 'नेपाली'),
('ng', 'Ndonga', 'Owambo'),
('nn', 'Norwegian Nynorsk', 'Norsk nynorsk'),
('no', 'Norwegian', 'Norsk'),
('ii', 'Nuosu', 'ꆈꌠ꒿ Nuosuhxop'),
('nr', 'South Ndebele', 'isiNdebele'),
('oc', 'Occitan', 'Occitan'),
('oj', 'Ojibwe', 'ᐊᓂᔑᓈᐯᒧᐎᓐ'),
('cu', 'Old Church Slavonic', 'ѩзыкъ словѣньскъ'),
('om', 'Oromo', 'Afaan Oromoo'),
('or', 'Oriya', 'ଓଡ଼ିଆ'),
('os', 'Ossetian', 'ирон æвзаг'),
('pa', 'Panjabi', 'ਪੰਜਾਬੀ'),
('pi', 'Pāli', 'पाऴि'),
('fa', 'Persian', 'فارسی'),
('pl', 'Polish', 'polski'),
('ps', 'Pashto', 'پښتو'),
('pt', 'Portuguese', 'Português'),
('qu', 'Quechua', 'Runa Simi'),
('rm', 'Romansh', 'rumantsch grischun'),
('rn', 'Kirundi', 'kiRundi'),
('ro', 'Romanian', 'română'),
('ru', 'Russian', 'русский язык'),
('sa', 'Sanskrit', 'संस्कृतम्'),
('sc', 'Sardinian', 'sardu'),
('sd', 'Sindhi', 'सिन्धी'),
('se', 'Northern Sami', 'Davvisámegiella'),
('sm', 'Samoan', 'gagana faa Samoa'),
('sg', 'Sango', 'yângâ tî sängö'),
('sr', 'Serbian', 'српски језик'),
('gd', 'Scottish Gaelic', 'Gàidhlig'),
('sn', 'Shona', 'chiShona'),
('si', 'Sinhala', 'සිංහල'),
('sk', 'Slovak', 'slovenčina'),
('sl', 'Slovene', 'slovenščina'),
('so', 'Somali', 'Soomaaliga'),
('st', 'Southern Sotho', 'Sesotho'),
('es', 'Spanish', 'español'),
('su', 'Sundanese', 'Basa Sunda'),
('sw', 'Swahili', 'Kiswahili'),
('ss', 'Swati', 'SiSwati'),
('sv', 'Swedish', 'svenska'),
('ta', 'Tamil', 'தமிழ்'),
('te', 'Telugu', 'తెలుగు'),
('tg', 'Tajik', 'тоҷикӣ'),
('th', 'Thai', 'ไทย'),
('ti', 'Tigrinya', 'ትግርኛ'),
('bo', 'Tibetan Standard', 'བོད་ཡིག'),
('tk', 'Turkmen', 'Türkmen'),
('tl', 'Tagalog', 'ᜏᜒᜃᜅ᜔ ᜆᜄᜎᜓᜄ᜔'),
('tn', 'Tswana', 'Setswana'),
('to', 'Tonga', 'faka Tonga'),
('tr', 'Turkish', 'Türkçe'),
('ts', 'Tsonga', 'Xitsonga'),
('tt', 'Tatar', 'татарча'),
('tw', 'Twi', 'Twi'),
('ty', 'Tahitian', 'Reo Tahiti'),
('ug', 'Uighur', 'ئۇيغۇرچە'),
('uk', 'Ukrainian', 'українська'),
('ur', 'Urdu', 'اردو'),
('uz', 'Uzbek', 'zbek'),
('ve', 'Venda', 'Tshivenḓa'),
('vi', 'Vietnamese', 'Tiếng Việt'),
('vo', 'Volapük', 'Volapük'),
('wa', 'Walloon', 'Walon'),
('cy', 'Welsh', 'Cymraeg'),
('wo', 'Wolof', 'Wollof'),
('fy', 'Western Frisian', 'Frysk'),
('xh', 'Xhosa', 'isiXhosa'),
('yi', 'Yiddish', 'ייִדיש'),
('yo', 'Yoruba', 'Yorùbá'),
('za', 'Zhuang', 'Saɯ cueŋƅ');

/* INITIALIZE OBJECT TYPES */
INSERT INTO object_types (id, name) VALUES
  (1, 'expression'),
  (2, 'text'),
  (3, 'article');

/* INITIALIZE USERS */
INSERT INTO users (id, email, username, encrypted_password, first_name, last_name, token, status, token_expires_at) VALUES
  (1, 'admin@akarpovich.com', 'admin', '$2a$14$1uz8bdnCERhrFJ1qDZ0gwOxxmHy4NuYsAu2mckpzL3r5C7WbO3nCO', '', '', '', 2, NOW()),
  (2, 'alexsure.k@gmail.com', 'akarpovich', '$2a$14$1uz8bdnCERhrFJ1qDZ0gwOxxmHy4NuYsAu2mckpzL3r5C7WbO3nCO', 'Aliaksandr', 'Karpovich', '', 2, NOW());


/* INITIALIZE GROUPS*/
INSERT INTO groups (id, target_lang_id, native_lang_id, name, status) VALUES 
  (1, 30, 134, 'Chinese', 0);

INSERT INTO user_group (user_id, group_id, role) VALUES
  (1, 1, 0),
  (2, 1, 0);

/* INITIALIZE SLICES */
INSERT INTO slices (id, name, visibility) VALUES
  (1, 'Slice 1', 1),
  (2, 'Slice 2', 1),
  (3, 'Slice 1.1', 1),
  (4, 'Slice 1.2', 1),
  (5, 'Slice 1.2.1', 1),
  (6, 'Slice 1.2.2', 1);

INSERT INTO group_slice (group_id, slice_id, path) VALUES
  (1, 1, ''),
  (1, 2, ''),
  (1, 3, '1'),
  (1, 4, '1'),
  (1, 5, '1.4'),
  (1, 6, '1.4');

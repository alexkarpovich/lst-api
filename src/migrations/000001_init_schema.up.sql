CREATE EXTENSION IF NOT EXISTS ltree;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TABLE IF EXISTS languages;
CREATE TABLE languages (
  code VARCHAR(2) PRIMARY KEY,
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

DROP TABLE IF EXISTS transcription_types;
CREATE TABLE transcription_types (
  id serial PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  lang VARCHAR(2) NOT NULL,
  CONSTRAINT fk_lang
    FOREIGN KEY(lang) 
    REFERENCES languages(code),
  UNIQUE(lang, name)
);

DROP TABLE IF EXISTS groups;
CREATE TABLE groups (
  id serial PRIMARY KEY,
  transcription_type INT NOT NULL,
  target_lang VARCHAR(2) NOT NULL,
  native_lang VARCHAR(2) NOT NULL,
  name VARCHAR(128) NOT NULL,
  config jsonb,
  status SMALLINT NOT NULL,
  CONSTRAINT fk_transcription_type
    FOREIGN KEY(transcription_type) 
    REFERENCES transcription_types(id),
  CONSTRAINT fk_target_lang
    FOREIGN KEY(target_lang) 
    REFERENCES languages(code),
  CONSTRAINT fk_native_lang
    FOREIGN KEY(native_lang) 
    REFERENCES languages(code),
  UNIQUE(name, target_lang, native_lang)
);

DROP TABLE IF EXISTS user_group;
CREATE TABLE user_group (
  user_id INT NOT NULL,
  group_id INT NOT NULL,
  role SMALLINT NOT NULL,
  token VARCHAR(128),
  token_expires_at TIMESTAMP,
  status SMALLINT NOT NULL,
  CONSTRAINT fk_user
    FOREIGN KEY(user_id) 
    REFERENCES users(id),
  CONSTRAINT fk_group
    FOREIGN KEY(group_id) 
    REFERENCES groups(id)
);

DROP TABLE IF EXISTS object_types;
CREATE TABLE object_types (
  id serial PRIMARY KEY,
  name VARCHAR(32) UNIQUE NOT NULL
);

DROP TABLE IF EXISTS expressions;
CREATE TABLE expressions (
  id serial PRIMARY KEY,
  lang VARCHAR(2) NOT NULL,
  value VARCHAR(128) NOT NULL,
  CONSTRAINT fk_lang
    FOREIGN KEY(lang) 
    REFERENCES languages(code),
  UNIQUE(lang, value)
);

DROP TABLE IF EXISTS translations;
CREATE TABLE translations (
  id serial PRIMARY KEY,
  target_id INT NOT NULL,
  native_id INT NOT NULL,
  type INT NOT NULL,
  comment VARCHAR(256),
  CONSTRAINT fk_type
    FOREIGN KEY(type) 
    REFERENCES object_types(id),
  UNIQUE(type, target_id, native_id)
);

DROP TABLE IF EXISTS texts;
CREATE TABLE texts (
  id serial PRIMARY KEY,
  origin_id INT,
  lang VARCHAR(2) NOT NULL,
  author_id INT NOT NULL,
  title VARCHAR(128) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_origin
    FOREIGN KEY(origin_id) 
    REFERENCES texts(id),
  CONSTRAINT fk_lang
    FOREIGN KEY(lang) 
    REFERENCES languages(code),
  CONSTRAINT fk_author
    FOREIGN KEY(author_id) 
    REFERENCES users(id)
);

DROP TABLE IF EXISTS transcriptions;
CREATE TABLE transcriptions (
  id serial PRIMARY KEY,
  type INT NOT NULL,
  value VARCHAR(128) NOT NULL,
  CONSTRAINT fk_type
    FOREIGN KEY(type) 
    REFERENCES transcription_types(id),
  UNIQUE(type, value)
);

DROP TABLE IF EXISTS expression_transcription;
CREATE TABLE expression_transcription (
  expression_id INT NOT NULL,
  transcription_id INT NOT NULL,
  CONSTRAINT fk_expression
    FOREIGN KEY(expression_id) 
    REFERENCES expressions(id),
  CONSTRAINT fk_transcription
    FOREIGN KEY(transcription_id) 
    REFERENCES transcriptions(id),
  UNIQUE(expression_id, transcription_id)
);

DROP TABLE IF EXISTS translation_transcription;
CREATE TABLE translation_transcription (
  translation_id INT NOT NULL,
  transcription_id INT NOT NULL,
  CONSTRAINT fk_translation
    FOREIGN KEY(translation_id) 
    REFERENCES translations(id),
  CONSTRAINT fk_transcription
    FOREIGN KEY(transcription_id) 
    REFERENCES transcriptions(id),
  UNIQUE(translation_id, transcription_id)
);

DROP TABLE IF EXISTS nodes;
CREATE TABLE nodes (
  id serial PRIMARY KEY,
  type SMALLINT NOT NULL,
  name VARCHAR(64) NOT NULL,
  visibility SMALLINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TRIGGER updated_at BEFORE UPDATE ON nodes
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


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


DROP TABLE IF EXISTS group_node;
CREATE TABLE group_node (
  group_id INT NOT NULL,
  node_id INT NOT NULL,
  path ltree NOT NULL,
  CONSTRAINT fk_group
    FOREIGN KEY(group_id) 
    REFERENCES groups(id),
  CONSTRAINT fk_node
    FOREIGN KEY(node_id) 
    REFERENCES nodes(id)
);

CREATE INDEX path_gist_idx ON group_node USING GIST (path);
CREATE INDEX path_idx ON group_node USING BTREE (path);

DROP TABLE IF EXISTS node_expression;
CREATE TABLE node_expression (
  node_id INT NOT NULL,
  expression_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_node
    FOREIGN KEY(node_id) 
    REFERENCES nodes(id),
  CONSTRAINT fk_expression
    FOREIGN KEY(expression_id) 
    REFERENCES expressions(id),
  UNIQUE(node_id, expression_id)
);

DROP TABLE IF EXISTS node_translation;
CREATE TABLE node_translation (
  node_id INT NOT NULL,
  translation_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_node
    FOREIGN KEY(node_id) 
    REFERENCES nodes(id),
  CONSTRAINT fk_translation
    FOREIGN KEY(translation_id) 
    REFERENCES translations(id),
  UNIQUE(node_id, translation_id)
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

DROP TABLE IF EXISTS trainings;
CREATE TABLE trainings (
  id serial PRIMARY KEY,
  owner_id INT NOT NULL,
  type SMALLINT NOT NULL,
  slices SMALLINT[] NOT NULL,
  UNIQUE(owner_id, type, slices)
);

DROP TABLE IF EXISTS training_items;
CREATE TABLE training_items (
  id serial PRIMARY KEY,
  training_id INT NOT NULL,
  expression_id INT NOT NULL,
  stage SMALLINT,
  cycle SMALLINT,
  complete BOOLEAN,
  CONSTRAINT fk_training
    FOREIGN KEY(training_id) 
    REFERENCES trainings(id),
  CONSTRAINT fk_expression
    FOREIGN KEY(expression_id) 
    REFERENCES expressions(id)
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

/* INITIALIZE DEFAULT TRANSCRIPTION TYPES */
INSERT INTO transcription_types (lang, name) VALUES
('ab', 'Default'),
('aa', 'Default'),
('af', 'Default'),
('ak', 'Default'),
('sq', 'Default'),
('am', 'Default'),
('ar', 'Default'),
('an', 'Default'),
('hy', 'Default'),
('as', 'Default'),
('av', 'Default'),
('ae', 'Default'),
('ay', 'Default'),
('az', 'Default'),
('bm', 'Default'),
('ba', 'Default'),
('eu', 'Default'),
('be', 'Default'),
('bn', 'Default'),
('bh', 'Default'),
('bi', 'Default'),
('bs', 'Default'),
('br', 'Default'),
('bg', 'Default'),
('my', 'Default'),
('ca', 'Default'),
('ch', 'Default'),
('ce', 'Default'),
('ny', 'Default'),
('zh', 'pinyin'),
('cv', 'Default'),
('kw', 'Default'),
('co', 'Default'),
('cr', 'Default'),
('hr', 'Default'),
('cs', 'Default'),
('da', 'Default'),
('dv', 'Default'),
('nl', 'Default'),
('en', 'Default'),
('eo', 'Default'),
('et', 'Default'),
('ee', 'Default'),
('fo', 'Default'),
('fj', 'Default'),
('fi', 'Default'),
('fr', 'Default'),
('ff', 'Default'),
('gl', 'Default'),
('ka', 'Default'),
('de', 'Default'),
('el', 'Default'),
('gn', 'Default'),
('gu', 'Default'),
('ht', 'Default'),
('ha', 'Default'),
('he', 'Default'),
('iw', 'Default'),
('hz', 'Default'),
('hi', 'Default'),
('ho', 'Default'),
('hu', 'Default'),
('ia', 'Default'),
('id', 'Default'),
('ie', 'Default'),
('ga', 'Default'),
('ig', 'Default'),
('ik', 'Default'),
('io', 'Default'),
('is', 'Default'),
('it', 'Default'),
('iu', 'Default'),
('ja', 'Default'),
('jv', 'Default'),
('kl', 'Default'),
('kn', 'Default'),
('kr', 'Default'),
('ks', 'Default'),
('kk', 'Default'),
('km', 'Default'),
('ki', 'Default'),
('rw', 'Default'),
('ky', 'Default'),
('kv', 'Default'),
('kg', 'Default'),
('ko', 'Default'),
('ku', 'Default'),
('kj', 'Default'),
('la', 'Default'),
('lb', 'Default'),
('lg', 'Default'),
('li', 'Default'),
('ln', 'Default'),
('lo', 'Default'),
('lt', 'Default'),
('lu', 'Default'),
('lv', 'Default'),
('gv', 'Default'),
('mk', 'Default'),
('mg', 'Default'),
('ms', 'Default'),
('ml', 'Default'),
('mt', 'Default'),
('mi', 'Default'),
('mr', 'Default'),
('mh', 'Default'),
('mn', 'Default'),
('na', 'Default'),
('nv', 'Default'),
('nb', 'Default'),
('nd', 'Default'),
('ne', 'Default'),
('ng', 'Default'),
('nn', 'Default'),
('no', 'Default'),
('ii', 'Default'),
('nr', 'Default'),
('oc', 'Default'),
('oj', 'Default'),
('cu', 'Default'),
('om', 'Default'),
('or', 'Default'),
('os', 'Default'),
('pa', 'Default'),
('pi', 'Default'),
('fa', 'Default'),
('pl', 'Default'),
('ps', 'Default'),
('pt', 'Default'),
('qu', 'Default'),
('rm', 'Default'),
('rn', 'Default'),
('ro', 'Default'),
('ru', 'Default'),
('sa', 'Default'),
('sc', 'Default'),
('sd', 'Default'),
('se', 'Default'),
('sm', 'Default'),
('sg', 'Default'),
('sr', 'Default'),
('gd', 'Default'),
('sn', 'Default'),
('si', 'Default'),
('sk', 'Default'),
('sl', 'Default'),
('so', 'Default'),
('st', 'Default'),
('es', 'Default'),
('su', 'Default'),
('sw', 'Default'),
('ss', 'Default'),
('sv', 'Default'),
('ta', 'Default'),
('te', 'Default'),
('tg', 'Default'),
('th', 'Default'),
('ti', 'Default'),
('bo', 'Default'),
('tk', 'Default'),
('tl', 'Default'),
('tn', 'Default'),
('to', 'Default'),
('tr', 'Default'),
('ts', 'Default'),
('tt', 'Default'),
('tw', 'Default'),
('ty', 'Default'),
('ug', 'Default'),
('uk', 'Default'),
('ur', 'Default'),
('uz', 'Default'),
('ve', 'Default'),
('vi', 'Default'),
('vo', 'Default'),
('wa', 'Default'),
('cy', 'Default'),
('wo', 'Default'),
('fy', 'Default'),
('xh', 'Default'),
('yi', 'Default'),
('yo', 'Default'),
('za', 'Default');

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
INSERT INTO groups (id, transcription_type, target_lang, native_lang, name, status) VALUES 
  (1, 30, 'zh', 'ru', 'Chinese', 0);

INSERT INTO user_group (user_id, group_id, role, status) VALUES
  (1, 1, 0, 1),
  (2, 1, 2, 1);

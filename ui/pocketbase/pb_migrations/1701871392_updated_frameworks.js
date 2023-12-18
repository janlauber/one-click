/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nwgspeyl1n5v30f")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "kv73rbuv",
    "name": "settings",
    "type": "json",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nwgspeyl1n5v30f")

  // remove
  collection.schema.removeField("kv73rbuv")

  return dao.saveCollection(collection)
})

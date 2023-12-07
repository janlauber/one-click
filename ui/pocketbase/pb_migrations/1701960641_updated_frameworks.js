/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nwgspeyl1n5v30f")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "y4kw2xt6",
    "name": "application",
    "type": "bool",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nwgspeyl1n5v30f")

  // remove
  collection.schema.removeField("y4kw2xt6")

  return dao.saveCollection(collection)
})
